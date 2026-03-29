package coursepass

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager struct {
	dbo  db.DB
	repo db.CoursesRepo
	auth AuthConfig
	embedlog.Logger
}

const (
	uxStudentsLoginConstraint = "uq_students_login"
	uxStudentsEmailConstraint = "uq_students_email"
	registerLockName          = "student_register"
)

func NewAuthManager(dbo db.DB, logger embedlog.Logger, authCfg AuthConfig) *AuthManager {
	return &AuthManager{
		dbo:    dbo,
		repo:   db.NewCoursesRepo(dbo),
		auth:   authCfg,
		Logger: logger,
	}
}

func (am *AuthManager) Register(ctx context.Context, login, password, email, firstName, lastName string) (*AuthToken, error) {
	hash, err := am.passwordHash(password)
	if err != nil {
		return nil, err
	}

	var authStudent *Student
	err = am.dbo.RunInLock(ctx, registerLockName, func(tx *pg.Tx) error {
		txRepo := am.repo.WithTransaction(tx)

		if err = am.ensureLoginAvailable(ctx, txRepo, login); err != nil {
			return err
		}
		if err = am.ensureEmailAvailable(ctx, txRepo, email); err != nil {
			return err
		}

		student, addErr := am.addStudent(ctx, txRepo, login, hash, firstName, lastName, email)
		if addErr != nil {
			return addErr
		}

		authStudent = student
		return nil
	})
	if err != nil {
		return nil, err
	}

	return am.newTokenForStudent(authStudent)
}

func (am *AuthManager) Login(ctx context.Context, login, password string) (*AuthToken, error) {
	student, err := am.studentByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	if err = am.checkPassword(student.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return am.newTokenForStudent(student)
}

func (am *AuthManager) passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed generate hash password: %w", err)
	}

	return string(hash), nil
}

func (am *AuthManager) isUniqueConstraintViolation(err error, constraintName string) bool {
	var pgErr pg.Error
	return errors.As(err, &pgErr) && pgErr.Field('n') == constraintName
}

func (am *AuthManager) ensureLoginAvailable(ctx context.Context, repo db.CoursesRepo, login string) error {
	studentData, err := repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return fmt.Errorf("failed check login: %w", err)
	}
	if studentData != nil {
		return ErrLoginExists
	}

	return nil
}

func (am *AuthManager) ensureEmailAvailable(ctx context.Context, repo db.CoursesRepo, email string) error {
	studentData, err := repo.OneStudent(ctx, &db.StudentSearch{
		Email: &email,
	})
	if err != nil {
		return fmt.Errorf("failed check email: %w", err)
	}
	if studentData != nil {
		return ErrEmailExists
	}

	return nil
}

func (am *AuthManager) addStudent(ctx context.Context, repo db.CoursesRepo, login, passwordHash, firstName, lastName, email string) (*Student, error) {
	student, err := repo.AddStudent(ctx, newDBStudent(login, passwordHash, firstName, lastName, email))
	if err != nil {
		if am.isUniqueConstraintViolation(err, uxStudentsLoginConstraint) {
			return nil, ErrLoginExists
		}
		if am.isUniqueConstraintViolation(err, uxStudentsEmailConstraint) {
			return nil, ErrEmailExists
		}

		return nil, fmt.Errorf("failed add student: %w", err)
	}

	domainStudent := newStudent(student)
	if domainStudent == nil {
		return nil, fmt.Errorf("failed add student: empty student")
	}

	return domainStudent, nil
}

func (am *AuthManager) studentByLogin(ctx context.Context, login string) (*Student, error) {
	studentData, err := am.repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return nil, ErrInvalidCredentials
	}

	domainStudent := newStudent(studentData)
	if domainStudent == nil {
		return nil, ErrInvalidCredentials
	}

	return domainStudent, nil
}

func (am *AuthManager) checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (am *AuthManager) newTokenForStudent(student *Student) (*AuthToken, error) {
	if student == nil {
		return nil, ErrInvalidCredentials
	}

	token, expiresIn, err := am.generateJWT(am.auth.JWTSecret, am.auth.JWTTTLSeconds, student.ID, student.Login)
	if err != nil {
		return nil, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

func (am *AuthManager) generateJWT(jwtSecretValue string, jwtTTLSeconds int, studentID int, login string) (string, int, error) {
	header := newTokenHeader()

	now := time.Now()
	ttl := am.tokenTTLSeconds(jwtTTLSeconds)
	exp := now.Unix() + int64(ttl)
	claims := newTokenClaims(studentID, login, now.Unix(), exp)

	headerRaw, err := am.marshalTokenHeader(header)
	if err != nil {
		return "", 0, err
	}
	claimsRaw, err := am.marshalTokenClaims(claims)
	if err != nil {
		return "", 0, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerRaw)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsRaw)
	unsigned := encodedHeader + "." + encodedClaims
	signature := am.signHS256(unsigned, am.jwtSecret(jwtSecretValue))

	return unsigned + "." + signature, ttl, nil
}

func (am *AuthManager) signHS256(unsignedToken string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsignedToken))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (am *AuthManager) marshalTokenClaims(c tokenClaims) ([]byte, error) {
	return json.Marshal(map[string]any{
		"sub":   c.Sub,
		"login": c.Login,
		"exp":   c.Exp,
		"iat":   c.Iat,
	})
}

func (am *AuthManager) marshalTokenHeader(h tokenHeader) ([]byte, error) {
	return json.Marshal(map[string]any{
		"alg": h.Alg,
		"typ": h.Typ,
	})
}

func (am *AuthManager) jwtSecret(secret string) []byte {
	if secret == "" {
		secret = "coursepass-dev-secret"
	}
	return []byte(secret)
}

func (am *AuthManager) tokenTTLSeconds(ttlSeconds int) int {
	if ttlSeconds <= 0 {
		return defaultTokenTTLSeconds
	}

	return ttlSeconds
}

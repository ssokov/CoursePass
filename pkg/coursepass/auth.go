package coursepass

import (
	"context"
	"errors"
	"fmt"

	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager struct {
	repo db.CoursesRepo
	auth AuthConfig
	embedlog.Logger
}

const (
	uxStudentsLoginConstraint = "uq_students_login"
	uxStudentsEmailConstraint = "uq_students_email"
)

func NewAuthManager(dbo db.DB, logger embedlog.Logger, authCfg AuthConfig) *AuthManager {
	return &AuthManager{
		repo:   db.NewCoursesRepo(dbo),
		auth:   authCfg,
		Logger: logger,
	}
}

func (am *AuthManager) Register(ctx context.Context, login, password, email, firstName, lastName string) (AuthToken, error) {
	if err := am.ensureLoginAvailable(ctx, login); err != nil {
		return AuthToken{}, err
	}
	if err := am.ensureEmailAvailable(ctx, email); err != nil {
		return AuthToken{}, err
	}

	hash, err := passwordHash(password)
	if err != nil {
		return AuthToken{}, err
	}

	student, err := am.addStudent(ctx, login, hash, firstName, lastName, email)
	if err != nil {
		return AuthToken{}, err
	}

	authStudent := newStudentAuth(*student)
	return am.newTokenForStudent(authStudent)
}

func (am *AuthManager) Login(ctx context.Context, login, password string) (AuthToken, error) {
	student, err := am.studentByLogin(ctx, login)
	if err != nil {
		return AuthToken{}, err
	}

	if err = checkPassword(student.PasswordHash, password); err != nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	return am.newTokenForStudent(student)
}

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed generate hash password: %w", err)
	}

	return string(hash), nil
}

func isUniqueConstraintViolation(err error, constraintName string) bool {
	var pgErr pg.Error
	return errors.As(err, &pgErr) && pgErr.Field('n') == constraintName
}

func (am *AuthManager) ensureLoginAvailable(ctx context.Context, login string) error {
	studentData, err := am.repo.OneStudent(ctx, &db.StudentSearch{
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

func (am *AuthManager) ensureEmailAvailable(ctx context.Context, email string) error {
	studentData, err := am.repo.OneStudent(ctx, &db.StudentSearch{
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

func (am *AuthManager) addStudent(ctx context.Context, login, passwordHash, firstName, lastName, email string) (*db.Student, error) {
	student, err := am.repo.AddStudent(ctx, newDBStudent(login, passwordHash, firstName, lastName, email))
	if err != nil {
		if isUniqueConstraintViolation(err, uxStudentsLoginConstraint) {
			return nil, ErrLoginExists
		}
		if isUniqueConstraintViolation(err, uxStudentsEmailConstraint) {
			return nil, ErrEmailExists
		}

		return nil, fmt.Errorf("failed add student: %w", err)
	}

	return student, nil
}

func (am *AuthManager) studentByLogin(ctx context.Context, login string) (studentAuth, error) {
	studentData, err := am.repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return studentAuth{}, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return studentAuth{}, ErrInvalidCredentials
	}

	return newStudentAuth(*studentData), nil
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (am *AuthManager) newTokenForStudent(student studentAuth) (AuthToken, error) {
	token, expiresIn, err := generateJWT(am.auth, student.StudentID, student.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

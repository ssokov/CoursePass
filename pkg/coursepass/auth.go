package coursepass

import (
	"context"
	"fmt"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager struct {
	repo db.CoursesRepo
	auth AuthConfig
	embedlog.Logger
}

func NewAuthManager(dbo db.DB, logger embedlog.Logger, authCfg AuthConfig) *AuthManager {
	return &AuthManager{
		repo:   db.NewCoursesRepo(dbo),
		auth:   authCfg,
		Logger: logger,
	}
}

func (am *AuthManager) Register(
	ctx context.Context,
	login, password, email, firstName, lastName string,
) (AuthToken, error) {
	if student, err := am.repo.OneStudent(ctx, &db.StudentSearch{Login: &login}); err != nil {
		return AuthToken{}, err
	} else if student != nil {
		return AuthToken{}, ErrLoginExists
	}

	if student, err := am.repo.OneStudent(ctx, &db.StudentSearch{Email: &email}); err != nil {
		return AuthToken{}, err
	} else if student != nil {
		return AuthToken{}, ErrEmailExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed generate hash password: %w", err)
	}

	student, err := am.repo.AddStudent(ctx, newDBStudent(login, string(passwordHash), firstName, lastName, email))
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create student: %w", err)
	}

	authStudent := newStudentAuth(*student)

	token, expiresIn, err := generateJWT(am.auth, authStudent.StudentID, authStudent.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

func (am *AuthManager) Login(ctx context.Context, login, password string) (AuthToken, error) {
	studentData, err := am.repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	student := newStudentAuth(*studentData)

	if err := bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(password)); err != nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	token, expiresIn, err := generateJWT(am.auth, student.StudentID, student.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

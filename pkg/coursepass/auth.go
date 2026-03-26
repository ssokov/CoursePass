package coursepass

import (
	"context"
	"fmt"

	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type AuthManager struct {
	db   db.DB
	repo db.CoursesRepo
	auth AuthConfig
	embedlog.Logger
}

const registerLockName = "student_register"

func NewAuthManager(dbo db.DB, logger embedlog.Logger, authCfg AuthConfig) *AuthManager {
	return &AuthManager{
		db:     dbo,
		repo:   db.NewCoursesRepo(dbo),
		auth:   authCfg,
		Logger: logger,
	}
}

func (am *AuthManager) Register(ctx context.Context, login, password, email, firstName, lastName string) (AuthToken, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed generate hash password: %w", err)
	}

	var student *db.Student
	err = am.db.RunInLock(ctx, registerLockName, func(tx *pg.Tx) error {
		txRepo := am.repo.WithTransaction(tx)

		if existing, checkErr := txRepo.OneStudent(ctx, &db.StudentSearch{Login: &login}); checkErr != nil {
			return checkErr
		} else if existing != nil {
			return ErrLoginExists
		}

		if existing, checkErr := txRepo.OneStudent(ctx, &db.StudentSearch{Email: &email}); checkErr != nil {
			return checkErr
		} else if existing != nil {
			return ErrEmailExists
		}

		var addErr error
		student, addErr = txRepo.AddStudent(ctx, newDBStudent(login, string(passwordHash), firstName, lastName, email))
		if addErr != nil {
			return fmt.Errorf("failed create student: %w", addErr)
		}

		return nil
	})
	if err != nil {
		return AuthToken{}, err
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

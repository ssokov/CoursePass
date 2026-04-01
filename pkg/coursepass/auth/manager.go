package auth

import (
	"context"
	"errors"
	"fmt"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	dbo            db.DB
	repo           db.CoursesRepo
	jwtSecretValue string
	jwtTTLSeconds  int
	embedlog.Logger
}

const (
	uxStudentsLoginConstraint = "uq_students_login"
	uxStudentsEmailConstraint = "uq_students_email"
	registerLockName          = "student_register"
)

func NewManager(dbo db.DB, logger embedlog.Logger, jwtSecretValue string, jwtTTLSeconds int) *Manager {
	return &Manager{
		dbo:            dbo,
		repo:           db.NewCoursesRepo(dbo),
		jwtSecretValue: jwtSecretValue,
		jwtTTLSeconds:  jwtTTLSeconds,
		Logger:         logger,
	}
}

func (am *Manager) validateModel(ctx context.Context, in any) error {
	var v coursepass.Validator
	v.CheckBasic(ctx, in)
	if v.HasErrors() {
		return v.Error()
	}
	return nil
}

func (am *Manager) ValidateStudent(ctx context.Context, draft coursepass.StudentDraft) error {
	return am.validateModel(ctx, draft)
}

func (am *Manager) ValidateStudentLogin(ctx context.Context, in coursepass.StudentLogin) error {
	return am.validateModel(ctx, in)
}

func (am *Manager) RegisterStudent(ctx context.Context, draft coursepass.StudentDraft) (*coursepass.AuthToken, error) {
	if err := am.ValidateStudent(ctx, draft); err != nil {
		return nil, coursepass.ErrValidation
	}

	hash, err := am.passwordHash(draft.Password)
	if err != nil {
		return nil, err
	}

	var authStudent *coursepass.Student
	err = am.dbo.RunInLock(ctx, registerLockName, func(tx *pg.Tx) error {
		txRepo := am.repo.WithTransaction(tx)

		if err = am.ensureLoginAvailable(ctx, txRepo, draft.Login); err != nil {
			return err
		}
		if err = am.ensureEmailAvailable(ctx, txRepo, draft.Email); err != nil {
			return err
		}

		student, addErr := am.addStudent(ctx, txRepo, draft.Login, hash, draft.FirstName, draft.LastName, draft.Email)
		if addErr != nil {
			return addErr
		}

		authStudent = student
		return nil
	})
	if err != nil {
		return nil, err
	}

	return am.NewTokenForStudent(authStudent)
}

func (am *Manager) Login(ctx context.Context, in coursepass.StudentLogin) (*coursepass.AuthToken, error) {
	if err := am.ValidateStudentLogin(ctx, in); err != nil {
		return nil, err
	}

	student, err := am.studentByLogin(ctx, in.Login)
	if err != nil {
		return nil, err
	}

	if err = am.checkPassword(student.PasswordHash, in.Password); err != nil {
		return nil, coursepass.ErrInvalidCredentials
	}

	return am.NewTokenForStudent(student)
}

func (am *Manager) passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed generate hash password: %w", err)
	}

	return string(hash), nil
}

func (am *Manager) isUniqueConstraintViolation(err error, constraintName string) bool {
	var pgErr pg.Error
	return errors.As(err, &pgErr) && pgErr.Field('n') == constraintName
}

func (am *Manager) ensureLoginAvailable(ctx context.Context, repo db.CoursesRepo, login string) error {
	studentData, err := repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return fmt.Errorf("failed check login: %w", err)
	}
	if studentData != nil {
		return coursepass.ErrLoginExists
	}

	return nil
}

func (am *Manager) ensureEmailAvailable(ctx context.Context, repo db.CoursesRepo, email string) error {
	studentData, err := repo.OneStudent(ctx, &db.StudentSearch{
		Email: &email,
	})
	if err != nil {
		return fmt.Errorf("failed check email: %w", err)
	}
	if studentData != nil {
		return coursepass.ErrEmailExists
	}

	return nil
}

func (am *Manager) addStudent(ctx context.Context, repo db.CoursesRepo, login, passwordHash, firstName, lastName, email string) (*coursepass.Student, error) {
	student, err := repo.AddStudent(ctx, coursepass.NewDBStudent(login, passwordHash, firstName, lastName, email))
	if err != nil {
		if am.isUniqueConstraintViolation(err, uxStudentsLoginConstraint) {
			return nil, coursepass.ErrLoginExists
		}
		if am.isUniqueConstraintViolation(err, uxStudentsEmailConstraint) {
			return nil, coursepass.ErrEmailExists
		}

		return nil, fmt.Errorf("failed add student: %w", err)
	}

	domainStudent := coursepass.NewStudent(student)
	if domainStudent == nil {
		return nil, fmt.Errorf("failed add student: empty student")
	}

	return domainStudent, nil
}

func (am *Manager) studentByLogin(ctx context.Context, login string) (*coursepass.Student, error) {
	studentData, err := am.repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return nil, coursepass.ErrInvalidCredentials
	}

	domainStudent := coursepass.NewStudent(studentData)
	if domainStudent == nil {
		return nil, coursepass.ErrInvalidCredentials
	}

	return domainStudent, nil
}

func (am *Manager) checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (am *Manager) NewTokenForStudent(student *coursepass.Student) (*coursepass.AuthToken, error) {
	if student == nil {
		return nil, coursepass.ErrInvalidCredentials
	}

	proc := newAuthProcessor(am.jwtSecretValue, am.jwtTTLSeconds)
	return proc.generateToken(student.ID, student.Login)
}

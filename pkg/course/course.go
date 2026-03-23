package course

import (
	"context"
	"fmt"
	"time"

	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
	"golang.org/x/crypto/bcrypt"
)

type CourseManager struct {
	db   db.DB
	repo db.CoursesRepo
	auth AuthConfig
	embedlog.Logger
}

func NewCourseManager(dbo db.DB, logger embedlog.Logger, authCfg AuthConfig) *CourseManager {
	return &CourseManager{
		db:     dbo,
		repo:   db.NewCoursesRepo(dbo),
		auth:   authCfg,
		Logger: logger,
	}
}

func (cm *CourseManager) Register(
	ctx context.Context,
	login, password, email, firstName, lastName string,
) (AuthToken, error) {
	if student, err := cm.repo.OneStudent(ctx, &db.StudentSearch{Login: &login}); err != nil {
		return AuthToken{}, err
	} else if student != nil {
		return AuthToken{}, ErrLoginExists
	}

	if student, err := cm.repo.OneStudent(ctx, &db.StudentSearch{Email: &email}); err != nil {
		return AuthToken{}, err
	} else if student != nil {
		return AuthToken{}, ErrEmailExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed generate hash password: %w", err)
	}

	student, err := cm.repo.AddStudent(ctx, newDBStudent(login, string(passwordHash), firstName, lastName, email))
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create student: %w", err)
	}

	token, expiresIn, err := generateJWT(cm.auth, student.ID, student.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

func (cm *CourseManager) Login(ctx context.Context, login, password string) (AuthToken, error) {
	student, err := cm.repo.OneStudent(ctx, &db.StudentSearch{
		Login: &login,
	})
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed get student: %w", err)
	}
	if student == nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(password)); err != nil {
		return AuthToken{}, ErrInvalidCredentials
	}

	token, expiresIn, err := generateJWT(cm.auth, student.ID, student.Login)
	if err != nil {
		return AuthToken{}, fmt.Errorf("failed create JWT: %w", err)
	}

	return newAuthToken(token, expiresIn), nil
}

func (cm *CourseManager) Me(ctx context.Context, studentID int) (*Student, error) {
	student, err := cm.repo.OneStudent(ctx, &db.StudentSearch{
		ID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if student == nil {
		return nil, ErrStudentNotFound
	}

	result := newStudent(*student)
	return &result, nil
}

func (cm *CourseManager) CoursesSummary(ctx context.Context, page, pageSize int) ([]CourseSummary, error) {
	currentTime := time.Now()

	courses, err := cm.repo.CoursesByFilters(ctx, &db.CourseSearch{
		AvailableFromTo: &currentTime,
		AvailableToFrom: &currentTime,
	}, db.Pager{
		Page:     page,
		PageSize: pageSize,
	})

	if err != nil {
		return nil, fmt.Errorf("failed get courses: %w", err)
	}

	return newCoursesSummary(courses), nil
}

func (cm *CourseManager) CourseByID(ctx context.Context, courseID int) (Course, error) {
	courseData, err := cm.repo.OneCourse(ctx, &db.CourseSearch{
		ID: &courseID,
	})
	if err != nil {
		return Course{}, fmt.Errorf("failed get course: %w", err)
	}
	if courseData == nil {
		return Course{}, ErrCourseNotFound
	}

	return newCourse(*courseData), nil
}

func (cm *CourseManager) StartExam(ctx context.Context, courseID, studentID int) (ExamStart, error) {
	if courseID <= 0 {
		return ExamStart{}, invalidCourseIDError()
	}

	currentTime := time.Now()
	var examStart ExamStart

	err := cm.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := cm.repo.WithTransaction(tx)

		courseData, err := txRepo.OneCourse(ctx, &db.CourseSearch{
			ID:              &courseID,
			AvailableFromTo: &currentTime,
			AvailableToFrom: &currentTime,
		})
		if err != nil {
			return fmt.Errorf("failed get course: %w", err)
		}
		if courseData == nil {
			return ErrCourseNotFound
		}

		questions, err := txRepo.QuestionsByFilters(
			ctx,
			&db.QuestionSearch{CourseID: &courseID},
			db.PagerNoLimit,
			db.WithSort(db.NewSortField(db.Columns.Question.ID, false)),
		)
		if err != nil {
			return fmt.Errorf("failed get questions: %w", err)
		}
		if len(questions) == 0 {
			return ErrNoQuestions
		}

		// TODO replace with colgen
		questionIDs := make([]int, len(questions))
		for i := range questions {
			questionIDs[i] = questions[i].ID
		}

		totalQuestions := len(questionIDs)
		examData, err := txRepo.AddExam(ctx, &db.Exam{
			CourseID:       courseID,
			StudentID:      studentID,
			QuestionIDs:    questionIDs,
			Answers:        db.ExamAnswers{},
			TotalQuestions: &totalQuestions,
			Status:         ExamStatusInProgress,
		})
		if err != nil {
			return fmt.Errorf("failed create exam: %w", err)
		}

		examStart = newExamStart(*examData, questionIDs)

		return nil
	})
	if err != nil {
		return ExamStart{}, err
	}

	return examStart, nil
}

func invalidCourseIDError() error {
	return ValidationError{
		Field:  "courseId",
		Reason: "must be greater than 0",
	}
}

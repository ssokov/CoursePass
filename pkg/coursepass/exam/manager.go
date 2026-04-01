package exam

import (
	"context"
	"errors"
	"fmt"
	"time"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
)

type Manager struct {
	db   db.DB
	repo db.CoursesRepo
	embedlog.Logger
}

const (
	ExamStatusPassed     = "passed"
	ExamStatusFailed     = "failed"
	ExamStatusInProgress = "in_progress"
	passScorePercent     = 70

	QuestionTypeSingleChoice   = "single_choice"
	QuestionTypeMultipleChoice = "multiple_choice"

	uxExamActiveStudentCourse = "ux_exams_active_student_course"
)

func NewManager(dbo db.DB, logger embedlog.Logger, _ string) *Manager {
	return &Manager{
		db:     dbo,
		repo:   db.NewCoursesRepo(dbo),
		Logger: logger,
	}
}

func (em *Manager) Start(ctx context.Context, studentID, courseID int) (*coursepass.Exam, error) {
	var exam *coursepass.Exam

	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)

		if err := em.getAvailableCourse(ctx, txRepo, courseID, time.Now()); err != nil {
			return err
		}

		questions, err := em.getCourseQuestions(ctx, txRepo, courseID)
		if err != nil {
			return err
		}

		examData, err := em.addExam(ctx, txRepo, courseID, studentID, questions.IDs())
		if err != nil {
			return err
		}

		exam = coursepass.NewExam(examData)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return exam, nil
}

func (em *Manager) Question(ctx context.Context, studentID, questionID, examID int) (*coursepass.Question, error) {
	exam, err := em.getExamByStudent(ctx, em.repo, studentID, examID)
	if err != nil {
		return nil, err
	}

	proc := newExamProcessor(*exam, nil)
	if err = proc.validateQuestionAccess(questionID); err != nil {
		return nil, err
	}

	return em.getQuestionForExam(ctx, em.repo, exam.CourseID, questionID)
}

func (em *Manager) SaveAnswer(ctx context.Context, studentID, examID, questionID int, optionIDs []int) error {
	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)

		exam, err := em.getExamByStudent(ctx, txRepo, studentID, examID)
		if err != nil {
			return err
		}

		question, err := em.getQuestionForAnswer(ctx, txRepo, exam.CourseID, questionID)
		if err != nil {
			return err
		}

		proc := newExamProcessor(*exam, nil)
		if err = proc.validateQuestionAccess(questionID); err != nil {
			return err
		}
		if err = proc.validateAnswer(*question, optionIDs); err != nil {
			return err
		}

		return em.updateExamAnswers(ctx, txRepo, exam.ID, proc.buildAnswers(questionID, optionIDs))
	})
	if err != nil {
		err = normalizeAnswerError(err)
		return fmt.Errorf("failed save answer: %w", err)
	}
	return nil
}

func (em *Manager) Submit(ctx context.Context, studentID, examID int) (*coursepass.Exam, error) {
	var result *coursepass.Exam

	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)

		exam, err := em.getExamByStudent(ctx, txRepo, studentID, examID)
		if err != nil {
			return err
		}

		questions, err := em.getQuestionsForSubmit(ctx, txRepo, exam.QuestionIDs)
		if err != nil {
			return err
		}

		proc := newExamProcessor(*exam, questions)
		if err = proc.validateSubmit(); err != nil {
			return err
		}
		sr := proc.calculateResult()

		finishedAt := time.Now()
		if err = em.updateSubmittedExam(ctx, txRepo, exam.ID, sr.status, sr.correctAnswers, sr.totalQuestions, sr.finalScore, finishedAt); err != nil {
			return err
		}

		finalScoreFloat := float64(sr.finalScore)
		exam.Status = sr.status
		exam.CorrectAnswers = &sr.correctAnswers
		exam.TotalQuestions = &sr.totalQuestions
		exam.FinalScore = &finalScoreFloat
		exam.FinishedAt = &finishedAt
		result = exam

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed submit exam: %w", err)
	}

	return result, nil
}

func (em *Manager) MyList(ctx context.Context, studentID, page, pageSize int) ([]coursepass.Exam, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return em.studentFinishedExams(ctx, studentID, page, pageSize)
}

func (em *Manager) getExamByStudent(ctx context.Context, repo db.CoursesRepo, studentID, examID int) (*coursepass.Exam, error) {
	examData, err := repo.OneExam(ctx, &db.ExamSearch{
		ID:        &examID,
		StudentID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get exam: %w", err)
	}
	if examData == nil {
		return nil, coursepass.ErrExamNotFound
	}

	return coursepass.NewExam(examData), nil
}

func (em *Manager) getQuestionForExam(ctx context.Context, repo db.CoursesRepo, courseID, questionID int) (*coursepass.Question, error) {
	questionData, err := repo.OneQuestion(
		ctx,
		&db.QuestionSearch{
			ID:       &questionID,
			CourseID: &courseID,
		},
		repo.FullQuestion(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed get question: %w", err)
	}
	if questionData == nil {
		return nil, coursepass.ErrQuestionNotFound
	}

	return coursepass.NewQuestion(questionData), nil
}

func (em *Manager) studentFinishedExams(ctx context.Context, studentID, page, pageSize int) ([]coursepass.Exam, error) {
	exams, err := em.repo.ExamsByFilters(ctx, &db.ExamSearch{
		StudentID: &studentID,
		StatusIn:  []string{ExamStatusPassed, ExamStatusFailed},
	}, db.Pager{
		Page:     page,
		PageSize: pageSize,
	},
		db.WithSort(db.NewSortField(db.Columns.Exam.FinishedAt, true)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed get exams: %w", err)
	}

	return coursepass.NewExams(exams), nil
}

func (em *Manager) getQuestionForAnswer(ctx context.Context, txRepo db.CoursesRepo, courseID, questionID int) (*coursepass.Question, error) {
	questionData, err := txRepo.OneQuestion(ctx, &db.QuestionSearch{
		ID:       &questionID,
		CourseID: &courseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get question: %w", err)
	}
	if questionData == nil {
		return nil, coursepass.ErrQuestionNotFound
	}

	return coursepass.NewQuestion(questionData), nil
}

func (em *Manager) getQuestionsForSubmit(ctx context.Context, txRepo db.CoursesRepo, questionIDs []int) ([]coursepass.Question, error) {
	questionData, err := txRepo.QuestionsByFilters(
		ctx,
		&db.QuestionSearch{IDs: questionIDs},
		db.PagerNoLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed get questions: %w", err)
	}

	return []coursepass.Question(coursepass.NewQuestions(questionData)), nil
}

func (em *Manager) getAvailableCourse(ctx context.Context, txRepo db.CoursesRepo, courseID int, now time.Time) error {
	courseData, err := txRepo.OneCourse(ctx, &db.CourseSearch{
		ID:              &courseID,
		AvailableFromTo: &now,
		AvailableToFrom: &now,
	})
	if err != nil {
		return fmt.Errorf("failed get coursepass: %w", err)
	}
	if courseData == nil {
		return coursepass.ErrCourseNotFound
	}

	return nil
}

func (em *Manager) getCourseQuestions(ctx context.Context, txRepo db.CoursesRepo, courseID int) (coursepass.Questions, error) {
	questions, err := txRepo.QuestionsByFilters(
		ctx,
		&db.QuestionSearch{CourseID: &courseID},
		db.PagerNoLimit,
		db.WithSort(db.NewSortField(db.Columns.Question.ID, false)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed get questions: %w", err)
	}
	if len(questions) == 0 {
		return nil, coursepass.ErrNoQuestions
	}

	return coursepass.NewQuestions(questions), nil
}

func (em *Manager) addExam(ctx context.Context, txRepo db.CoursesRepo, courseID, studentID int, questionIDs []int) (*db.Exam, error) {
	totalQuestions := len(questionIDs)
	examData, err := txRepo.AddExam(ctx, &db.Exam{
		CourseID:       courseID,
		StudentID:      studentID,
		QuestionIDs:    questionIDs,
		Answers:        db.ExamAnswers{},
		TotalQuestions: &totalQuestions,
		Status:         ExamStatusInProgress,
	})
	if err != nil && isActiveExamUniqueViolation(err) {
		return nil, coursepass.ErrExamAlreadyStarted
	}
	if err != nil {
		return nil, fmt.Errorf("failed add exam: %w", err)
	}

	return examData, nil
}

func (em *Manager) updateExamAnswers(ctx context.Context, txRepo db.CoursesRepo, examID int, answers db.ExamAnswers) error {
	updated, err := txRepo.UpdateExam(
		ctx,
		coursepass.NewDBExamAnswersUpdate(examID, answers),
		db.WithColumns(db.Columns.Exam.Answers),
	)
	if err != nil {
		return fmt.Errorf("failed update exam: %w", err)
	}
	if !updated {
		return coursepass.ErrExamNotUpdated
	}

	return nil
}

func (em *Manager) updateSubmittedExam(ctx context.Context, txRepo db.CoursesRepo, examID int, status string, correctAnswers, totalQuestions, finalScore int, finishedAt time.Time) error {
	updated, err := txRepo.UpdateExam(
		ctx,
		coursepass.NewDBExamSubmitUpdate(examID, status, correctAnswers, totalQuestions, float64(finalScore), finishedAt),
		db.WithColumns(
			db.Columns.Exam.Status,
			db.Columns.Exam.CorrectAnswers,
			db.Columns.Exam.TotalQuestions,
			db.Columns.Exam.FinalScore,
			db.Columns.Exam.FinishedAt,
		),
	)
	if err != nil {
		return fmt.Errorf("failed update exam: %w", err)
	}
	if !updated {
		return coursepass.ErrExamNotUpdated
	}

	return nil
}

func isActiveExamUniqueViolation(err error) bool {
	var pgErr pg.Error
	return errors.As(err, &pgErr) && pgErr.Field('n') == uxExamActiveStudentCourse
}

func normalizeAnswerError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, coursepass.ErrExamNotInProgress),
		errors.Is(err, coursepass.ErrQuestionNotInExam),
		errors.Is(err, coursepass.ErrAnswerAlreadySaved),
		errors.Is(err, coursepass.ErrQuestionNotFound):
		return coursepass.ErrAnswerUnavailable
	default:
		return err
	}
}

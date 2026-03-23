package coursepass

import (
	"context"
	"fmt"
	"slices"
	"time"

	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
)

type ExamManager struct {
	db           db.DB
	repo         db.CoursesRepo
	mediaWebPath string
	embedlog.Logger
}

func NewExamManager(dbo db.DB, logger embedlog.Logger, mediaWebPath string) *ExamManager {
	return &ExamManager{
		db:           dbo,
		repo:         db.NewCoursesRepo(dbo),
		mediaWebPath: mediaWebPath,
		Logger:       logger,
	}
}

func (em *ExamManager) Start(ctx context.Context, courseID, studentID int) (ExamStart, error) {
	if courseID <= 0 {
		return ExamStart{}, newValidationError("courseId", "must be greater than 0")
	}

	currentTime := time.Now()
	var examStart ExamStart

	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)

		courseData, err := txRepo.OneCourse(ctx, &db.CourseSearch{
			ID:              &courseID,
			AvailableFromTo: &currentTime,
			AvailableToFrom: &currentTime,
		})
		if err != nil {
			return fmt.Errorf("failed get coursepass: %w", err)
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

func (em *ExamManager) Question(ctx context.Context, examID, questionID, studentID int) (Question, error) {
	examData, err := em.repo.OneExam(ctx, &db.ExamSearch{
		ID:        &examID,
		StudentID: &studentID,
	})
	if err != nil {
		return Question{}, fmt.Errorf("failed get exam: %w", err)
	}
	if examData == nil {
		return Question{}, ErrExamNotFound
	}
	if examData.Status != ExamStatusInProgress {
		return Question{}, ErrExamNotInProgress
	}
	if !slices.Contains(examData.QuestionIDs, questionID) {
		return Question{}, newValidationError("questionId", "does not belong to exam")
	}

	questionData, err := em.repo.OneQuestion(ctx, &db.QuestionSearch{
		ID:       &questionID,
		CourseID: &examData.CourseID,
	}, em.repo.FullQuestion())
	if err != nil {
		return Question{}, fmt.Errorf("failed get question: %w", err)
	}
	if questionData == nil {
		return Question{}, ErrQuestionNotFound
	}

	return newQuestion(*questionData, em.mediaWebPath), nil
}

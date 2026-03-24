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

		// TODO replace with colgen (не очень понимаю, как это сделать)
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
		return Question{}, ErrQuestionNotInExam
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

func (em *ExamManager) SaveAnswer(ctx context.Context, studentID, examID, questionID int, optionIDs []int) error {
	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)
		exam, err := txRepo.OneExam(ctx, &db.ExamSearch{
			ID:        &examID,
			StudentID: &studentID,
		})
		if err != nil {
			return fmt.Errorf("failed get exam: %w", err)
		}
		if exam == nil {
			return ErrExamNotFound
		}

		if exam.Status != ExamStatusInProgress {
			return ErrExamNotInProgress
		}

		// The question must belong to the exam snapshot captured at Start.
		if !slices.Contains(exam.QuestionIDs, questionID) {
			return ErrQuestionNotInExam
		}

		// A question can be answered only once.
		for _, ans := range exam.Answers {
			if ans.QuestionID == questionID {
				return ErrAnswerAlreadySaved
			}
		}

		question, err := txRepo.OneQuestion(ctx, &db.QuestionSearch{
			ID:       &questionID,
			CourseID: &exam.CourseID,
		})

		if err != nil {
			return fmt.Errorf("failed get question: %w", err)
		}
		if question == nil {
			return ErrQuestionNotFound
		}

		// Validate that every optionID belongs to this question.
		allowed := make(map[int]struct{}, len(question.Options))
		for _, opt := range question.Options {
			allowed[opt.OptionID] = struct{}{}
		}
		for _, id := range optionIDs {
			if _, ok := allowed[id]; !ok {
				return ErrInvalidOptionIDs
			}
		}

		if len(optionIDs) > 1 && question.QuestionType == QuestionTypeSingleChoice {
			return ErrInvalidOptionIDs
		}

		exam.Answers = append(exam.Answers, db.ExamAnswer{QuestionID: questionID, OptionIDs: optionIDs})
		flag, err := txRepo.UpdateExam(ctx, exam, db.WithColumns(db.Columns.Exam.Answers))

		if err != nil {
			return fmt.Errorf("failed update exam: %w", err)
		}
		if !flag {
			return ErrExamNotUpdated
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed save answer: %w", err)
	}
	return nil
}

package coursepass

import (
	"context"
	"fmt"
	"math"
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

const (
	ExamStatusPassed     = "passed"
	ExamStatusFailed     = "failed"
	ExamStatusInProgress = "in_progress"
	passScorePercent     = 70

	QuestionTypeSingleChoice   = "single_choice"
	QuestionTypeMultipleChoice = "multiple_choice"
)

func NewExamManager(dbo db.DB, logger embedlog.Logger, mediaWebPath string) *ExamManager {
	return &ExamManager{
		db:           dbo,
		repo:         db.NewCoursesRepo(dbo),
		mediaWebPath: mediaWebPath,
		Logger:       logger,
	}
}

func (em *ExamManager) Start(ctx context.Context, studentID, courseID int) (ExamStart, error) {
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

		courseQuestions := newQuestions(questions, em.mediaWebPath)
		questionIDs := Questions(courseQuestions).QuestionIDs()

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

func (em *ExamManager) Question(ctx context.Context, studentID, questionID, examID int) (Question, error) {
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

	exam := newExamState(*examData)
	if exam.Status != ExamStatusInProgress {
		return Question{}, ErrExamNotInProgress
	}
	if !slices.Contains(exam.QuestionIDs, questionID) {
		return Question{}, ErrQuestionNotInExam
	}

	questionData, err := em.repo.OneQuestion(ctx, &db.QuestionSearch{
		ID:       &questionID,
		CourseID: &exam.CourseID,
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
		examData, err := txRepo.OneExam(ctx, &db.ExamSearch{
			ID:        &examID,
			StudentID: &studentID,
		})
		if err != nil {
			return fmt.Errorf("failed get exam: %w", err)
		}
		if examData == nil {
			return ErrExamNotFound
		}

		exam := newExamState(*examData)

		if exam.Status != ExamStatusInProgress {
			return ErrExamNotInProgress
		}

		// The question must belong to the exam snapshot captured at Start.
		if !slices.Contains(exam.QuestionIDs, questionID) {
			return ErrQuestionNotInExam
		}

		answerByQuestionID := ExamAnswers(exam.Answers).IndexByQuestionID()
		if _, exists := answerByQuestionID[questionID]; exists {
			return ErrAnswerAlreadySaved
		}

		questionData, err := txRepo.OneQuestion(ctx, &db.QuestionSearch{
			ID:       &questionID,
			CourseID: &exam.CourseID,
		})

		if err != nil {
			return fmt.Errorf("failed get question: %w", err)
		}
		if questionData == nil {
			return ErrQuestionNotFound
		}

		question := newQuestion(*questionData, em.mediaWebPath)

		allowedOptionByID := QuestionOptions(question.Options).IndexByOptionID()
		for _, id := range optionIDs {
			if _, ok := allowedOptionByID[id]; !ok {
				return ErrInvalidOptionIDs
			}
		}

		if len(optionIDs) > 1 && question.QuestionType == QuestionTypeSingleChoice {
			return ErrInvalidOptionIDs
		}

		exam.Answers = append(exam.Answers, ExamAnswer{
			QuestionID: questionID,
			OptionIDs:  slicesClone(optionIDs),
		})

		flag, err := txRepo.UpdateExam(
			ctx,
			newDBExamAnswersUpdate(exam.ExamID, exam.Answers),
			db.WithColumns(db.Columns.Exam.Answers),
		)

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

func (em *ExamManager) Submit(ctx context.Context, studentID, examID int) (ExamResult, error) {
	var examResult ExamResult

	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)
		examData, err := txRepo.OneExam(ctx, &db.ExamSearch{
			ID:        &examID,
			StudentID: &studentID,
		})
		if err != nil {
			return fmt.Errorf("failed get exam: %w", err)
		}
		if examData == nil {
			return ErrExamNotFound
		}
		exam := newExamState(*examData)
		if exam.Status != ExamStatusInProgress {
			return ErrExamNotInProgress
		}

		questionIDs := exam.QuestionIDs
		totalQuestions := len(questionIDs)
		if totalQuestions == 0 {
			return ErrNoQuestions
		}

		questionData, err := txRepo.QuestionsByFilters(
			ctx,
			&db.QuestionSearch{IDs: questionIDs},
			db.PagerNoLimit,
		)
		if err != nil {
			return fmt.Errorf("failed get questions: %w", err)
		}

		questions := newQuestions(questionData, em.mediaWebPath)
		questionByID := Questions(questions).IndexByQuestionID()
		answerByQuestionID := ExamAnswers(exam.Answers).IndexByQuestionID()

		var correctAnswers int
		for _, questionID := range questionIDs {
			question, ok := questionByID[questionID]
			if !ok {
				continue
			}

			correctOptionIDs := getCorrectOptionIDs(question.Options)
			answer, hasAnswer := answerByQuestionID[questionID]
			if !hasAnswer {
				continue
			}

			if equalOptionIDSets(correctOptionIDs, answer.OptionIDs) {
				correctAnswers++
			}
		}

		finalScore := calculateFinalScore(correctAnswers, totalQuestions)
		status := ExamStatusFailed
		if finalScore >= passScorePercent {
			status = ExamStatusPassed
		}

		finishedAt := time.Now()
		finalScoreFloat := float64(finalScore)

		updated, err := txRepo.UpdateExam(
			ctx,
			newDBExamSubmitUpdate(
				exam.ExamID,
				status,
				correctAnswers,
				totalQuestions,
				finalScoreFloat,
				finishedAt,
			),
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
			return ErrExamNotUpdated
		}

		examResult = ExamResult{
			ExamID:         exam.ExamID,
			Status:         status,
			FinalScore:     finalScore,
			CorrectAnswers: correctAnswers,
			TotalQuestions: totalQuestions,
		}

		return nil
	})
	if err != nil {
		return ExamResult{}, fmt.Errorf("failed submit exam: %w", err)
	}

	return examResult, nil
}

func (em *ExamManager) MyList(ctx context.Context, studentID, page, pageSize int) ([]ExamSummary, error) {

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

	return newExamSummaries(exams), nil
}

func getCorrectOptionIDs(options []QuestionOption) []int {
	optionByCorrectness := QuestionOptions(options).GroupByIsCorrect()
	correctOptions, ok := optionByCorrectness[true]
	if !ok {
		return nil
	}

	return QuestionOptions(correctOptions).OptionIDs()
}

func equalOptionIDSets(a, b []int) bool {
	setA := make(map[int]struct{}, len(a))
	for _, id := range a {
		setA[id] = struct{}{}
	}

	setB := make(map[int]struct{}, len(b))
	for _, id := range b {
		setB[id] = struct{}{}
	}

	if len(setA) != len(setB) {
		return false
	}

	for id := range setA {
		if _, ok := setB[id]; !ok {
			return false
		}
	}

	return true
}

func calculateFinalScore(correctAnswers, totalQuestions int) int {
	if totalQuestions <= 0 {
		return 0
	}

	score := (float64(correctAnswers) * 100) / float64(totalQuestions)
	return int(math.Round(score))
}

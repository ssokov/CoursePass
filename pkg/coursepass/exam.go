package coursepass

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"time"

	"courses/pkg/db"

	"github.com/go-pg/pg/v10"
	"github.com/vmkteam/embedlog"
)

type ExamManager struct {
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

func NewExamManager(dbo db.DB, logger embedlog.Logger, _ string) *ExamManager {
	return &ExamManager{
		db:     dbo,
		repo:   db.NewCoursesRepo(dbo),
		Logger: logger,
	}
}

func (em *ExamManager) Start(ctx context.Context, studentID, courseID int) (*Exam, error) {
	var exam *Exam

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

		exam = NewExam(examData)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return exam, nil
}

func (em *ExamManager) Question(ctx context.Context, studentID, questionID, examID int) (*Question, error) {
	exam, err := em.getExamByStudent(ctx, em.repo, studentID, examID)
	if err != nil {
		return nil, err
	}
	if exam.Status != ExamStatusInProgress {
		return nil, ErrExamNotInProgress
	}
	if !slices.Contains(exam.QuestionIDs, questionID) {
		return nil, ErrQuestionNotInExam
	}

	return em.getQuestionForExam(ctx, em.repo, exam.CourseID, questionID)
}

func (em *ExamManager) SaveAnswer(ctx context.Context, studentID, examID, questionID int, optionIDs []int) error {
	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)
		exam, err := em.getExamForAnswer(ctx, txRepo, studentID, examID, questionID)
		if err != nil {
			return err
		}

		question, err := em.getQuestionForAnswer(ctx, txRepo, exam.CourseID, questionID)
		if err != nil {
			return err
		}

		if err = validateAnswerOptions(*question, optionIDs); err != nil {
			return err
		}

		exam.Answers = append(exam.Answers, db.ExamAnswer{
			QuestionID: questionID,
			OptionIDs:  slices.Clone(optionIDs),
		})
		return em.updateExamAnswers(ctx, txRepo, exam.ID, exam.Answers)
	})
	if err != nil {
		return fmt.Errorf("failed save answer: %w", err)
	}
	return nil
}

func (em *ExamManager) Submit(ctx context.Context, studentID, examID int) (*Exam, error) {
	var result *Exam

	err := em.db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		txRepo := em.repo.WithTransaction(tx)
		exam, err := em.getExamForSubmit(ctx, txRepo, studentID, examID)
		if err != nil {
			return err
		}

		questions, err := em.getQuestionsForSubmit(ctx, txRepo, exam.QuestionIDs)
		if err != nil {
			return err
		}

		status, correctAnswers, totalQuestions, finalScore := calculateSubmitMetrics(*exam, questions)
		finishedAt := time.Now()
		if err = em.updateSubmittedExam(ctx, txRepo, exam.ID, status, correctAnswers, totalQuestions, finalScore, finishedAt); err != nil {
			return err
		}

		finalScoreFloat := float64(finalScore)
		exam.Status = status
		exam.CorrectAnswers = &correctAnswers
		exam.TotalQuestions = &totalQuestions
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

func (em *ExamManager) MyList(ctx context.Context, studentID, page, pageSize int) ([]Exam, error) {
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

	return NewExams(exams), nil
}

func (em *ExamManager) getExamByStudent(ctx context.Context, repo db.CoursesRepo, studentID, examID int) (*Exam, error) {
	examData, err := repo.OneExam(ctx, &db.ExamSearch{
		ID:        &examID,
		StudentID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get exam: %w", err)
	}
	if examData == nil {
		return nil, ErrExamNotFound
	}

	return NewExam(examData), nil
}

func (em *ExamManager) getQuestionForExam(ctx context.Context, repo db.CoursesRepo, courseID, questionID int) (*Question, error) {
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
		return nil, ErrQuestionNotFound
	}

	return NewQuestion(questionData), nil
}

func (em *ExamManager) getExamForAnswer(ctx context.Context, txRepo db.CoursesRepo, studentID, examID, questionID int) (*Exam, error) {
	examData, err := txRepo.OneExam(ctx, &db.ExamSearch{
		ID:        &examID,
		StudentID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get exam: %w", err)
	}
	if examData == nil {
		return nil, ErrExamNotFound
	}

	exam := NewExam(examData)
	if err = validateExamQuestionAccess(*exam, questionID); err != nil {
		return nil, err
	}

	answerByQuestionID := ExamAnswers(exam.Answers).IndexByQuestionID()
	if _, exists := answerByQuestionID[questionID]; exists {
		return nil, ErrAnswerAlreadySaved
	}

	return exam, nil
}

func (em *ExamManager) getQuestionForAnswer(ctx context.Context, txRepo db.CoursesRepo, courseID, questionID int) (*Question, error) {
	questionData, err := txRepo.OneQuestion(ctx, &db.QuestionSearch{
		ID:       &questionID,
		CourseID: &courseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get question: %w", err)
	}
	if questionData == nil {
		return nil, ErrQuestionNotFound
	}

	return NewQuestion(questionData), nil
}

func (em *ExamManager) getExamForSubmit(ctx context.Context, txRepo db.CoursesRepo, studentID, examID int) (*Exam, error) {
	examData, err := txRepo.OneExam(ctx, &db.ExamSearch{
		ID:        &examID,
		StudentID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get exam: %w", err)
	}
	if examData == nil {
		return nil, ErrExamNotFound
	}

	exam := NewExam(examData)
	if exam.Status != ExamStatusInProgress {
		return nil, ErrExamNotInProgress
	}

	return exam, nil
}

func (em *ExamManager) getQuestionsForSubmit(ctx context.Context, txRepo db.CoursesRepo, questionIDs []int) ([]Question, error) {
	if len(questionIDs) == 0 {
		return nil, ErrNoQuestions
	}

	questionData, err := txRepo.QuestionsByFilters(
		ctx,
		&db.QuestionSearch{IDs: questionIDs},
		db.PagerNoLimit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed get questions: %w", err)
	}

	return []Question(NewQuestions(questionData)), nil
}

func (em *ExamManager) getAvailableCourse(ctx context.Context, txRepo db.CoursesRepo, courseID int, now time.Time) error {
	courseData, err := txRepo.OneCourse(ctx, &db.CourseSearch{
		ID:              &courseID,
		AvailableFromTo: &now,
		AvailableToFrom: &now,
	})
	if err != nil {
		return fmt.Errorf("failed get coursepass: %w", err)
	}
	if courseData == nil {
		return ErrCourseNotFound
	}

	return nil
}

func (em *ExamManager) getCourseQuestions(ctx context.Context, txRepo db.CoursesRepo, courseID int) (Questions, error) {
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
		return nil, ErrNoQuestions
	}

	return NewQuestions(questions), nil
}

func (em *ExamManager) addExam(ctx context.Context, txRepo db.CoursesRepo, courseID, studentID int, questionIDs []int) (*db.Exam, error) {
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
		return nil, ErrExamAlreadyStarted
	}
	if err != nil {
		return nil, fmt.Errorf("failed add exam: %w", err)
	}

	return examData, nil
}

func (em *ExamManager) updateExamAnswers(ctx context.Context, txRepo db.CoursesRepo, examID int, answers db.ExamAnswers) error {
	updated, err := txRepo.UpdateExam(
		ctx,
		newDBExamAnswersUpdate(examID, answers),
		db.WithColumns(db.Columns.Exam.Answers),
	)
	if err != nil {
		return fmt.Errorf("failed update exam: %w", err)
	}
	if !updated {
		return ErrExamNotUpdated
	}

	return nil
}

func (em *ExamManager) updateSubmittedExam(ctx context.Context, txRepo db.CoursesRepo, examID int, status string, correctAnswers, totalQuestions, finalScore int, finishedAt time.Time) error {
	updated, err := txRepo.UpdateExam(
		ctx,
		newDBExamSubmitUpdate(examID, status, correctAnswers, totalQuestions, float64(finalScore), finishedAt),
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

	return nil
}

func validateExamQuestionAccess(exam Exam, questionID int) error {
	if exam.Status != ExamStatusInProgress {
		return ErrExamNotInProgress
	}
	if !slices.Contains(exam.QuestionIDs, questionID) {
		return ErrQuestionNotInExam
	}

	return nil
}

func validateAnswerOptions(question Question, optionIDs []int) error {
	allowedOptionByID := QuestionOptions(question.Options).IndexByOptionID()
	for _, id := range optionIDs {
		if _, ok := allowedOptionByID[id]; !ok {
			return ErrInvalidOptionIDs
		}
	}

	if len(optionIDs) > 1 && question.QuestionType == QuestionTypeSingleChoice {
		return ErrInvalidOptionIDs
	}

	return nil
}

func calculateSubmitMetrics(exam Exam, questions []Question) (string, int, int, int) {
	totalQuestions := len(exam.QuestionIDs)
	correctAnswers := countCorrectAnswers(exam.QuestionIDs, questions, ExamAnswers(exam.Answers))
	finalScore := calculateFinalScore(correctAnswers, totalQuestions)

	status := ExamStatusFailed
	if finalScore >= passScorePercent {
		status = ExamStatusPassed
	}

	return status, correctAnswers, totalQuestions, finalScore
}

func isActiveExamUniqueViolation(err error) bool {
	var pgErr pg.Error
	return errors.As(err, &pgErr) && pgErr.Field('n') == uxExamActiveStudentCourse
}

func countCorrectAnswers(questionIDs []int, questions Questions, answers ExamAnswers) int {
	questionByID := questions.Index()
	answerByQuestionID := answers.IndexByQuestionID()

	var correctAnswers int
	for _, questionID := range questionIDs {
		question, ok := questionByID[questionID]
		if !ok {
			continue
		}

		correctOptionIDs := getCorrectOptionIDs(QuestionOptions(question.Options))
		answer, hasAnswer := answerByQuestionID[questionID]
		if !hasAnswer {
			continue
		}

		if equalOptionIDSets(correctOptionIDs, answer.OptionIDs) {
			correctAnswers++
		}
	}

	return correctAnswers
}

func getCorrectOptionIDs(options QuestionOptions) []int {
	optionByCorrectness := options.GroupByIsCorrect()
	correctOptions, ok := optionByCorrectness[true]
	if !ok {
		return nil
	}

	return correctOptions.OptionIDs()
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

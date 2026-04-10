package exam

import (
	"math"
	"slices"

	"courses/pkg/coursepass"
	"courses/pkg/db"
)

type processor struct {
	exam      coursepass.Exam
	questions []coursepass.Question

	// ExamManager
}

func newExamProcessor(exam coursepass.Exam, questions []coursepass.Question) *processor {
	return &processor{exam: exam, questions: questions}
}

type submitResult struct {
	status         string
	correctAnswers int
	totalQuestions int
	finalScore     int
}

func (p *processor) validateQuestionAccess(questionID int) error {
	if p.exam.Status != ExamStatusInProgress {
		return coursepass.ErrExamNotInProgress
	}
	if !slices.Contains(p.exam.QuestionIDs, questionID) {
		return coursepass.ErrQuestionNotInExam
	}

	return nil
}

func (p *processor) validateAnswer(question coursepass.Question, optionIDs []int) error {
	answerByQuestionID := coursepass.ExamAnswers(p.exam.Answers).IndexByQuestionID()
	if _, exists := answerByQuestionID[question.ID]; exists {
		return coursepass.ErrAnswerAlreadySaved
	}

	allowedOptionByID := coursepass.QuestionOptions(question.Options).IndexByOptionID()
	for _, id := range optionIDs {
		if _, ok := allowedOptionByID[id]; !ok {
			return coursepass.ErrInvalidOptionIDs
		}
	}

	if len(optionIDs) > 1 && question.QuestionType == QuestionTypeSingleChoice {
		return coursepass.ErrInvalidOptionIDs
	}

	return nil
}

func (p *processor) buildAnswers(questionID int, optionIDs []int) db.ExamAnswers {
	return append(p.exam.Answers, db.ExamAnswer{
		QuestionID: questionID,
		OptionIDs:  slices.Clone(optionIDs),
	})
}

func (p *processor) validateSubmit() error {
	if p.exam.Status != ExamStatusInProgress {
		return coursepass.ErrExamNotInProgress
	}
	if len(p.exam.QuestionIDs) == 0 {
		return coursepass.ErrNoQuestions
	}

	return nil
}

func (p *processor) calculateResult() submitResult {
	totalQuestions := len(p.exam.QuestionIDs)
	correctAnswers := p.countCorrectAnswers()
	finalScore := calculateFinalScore(correctAnswers, totalQuestions)

	status := ExamStatusFailed
	if finalScore >= passScorePercent {
		status = ExamStatusPassed
	}

	return submitResult{
		status:         status,
		correctAnswers: correctAnswers,
		totalQuestions: totalQuestions,
		finalScore:     finalScore,
	}
}

func (p *processor) countCorrectAnswers() int {
	questionByID := coursepass.Questions(p.questions).Index()
	answerByQuestionID := coursepass.ExamAnswers(p.exam.Answers).IndexByQuestionID()

	var correctAnswers int
	for _, questionID := range p.exam.QuestionIDs {
		question, ok := questionByID[questionID]
		if !ok {
			continue
		}

		correctOptionIDs := getCorrectOptionIDs(coursepass.QuestionOptions(question.Options))
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

func getCorrectOptionIDs(options coursepass.QuestionOptions) []int {
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

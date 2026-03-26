package coursepass

import (
	"errors"
	"fmt"
)

var (
	ErrValidation         = errors.New("validation error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrLoginExists        = errors.New("login already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrStudentNotFound    = errors.New("student not found")
	ErrCourseNotFound     = errors.New("coursepass not found")
	ErrExamNotFound       = errors.New("exam not found")
	ErrExamNotInProgress  = errors.New("exam is not in progress")
	ErrQuestionNotFound   = errors.New("question not found")
	ErrQuestionNotInExam  = errors.New("question does not belong to exam")
	ErrNoQuestions        = errors.New("coursepass has no questions")
	ErrAnswerAlreadySaved = errors.New("answer already saved")
	ErrInvalidOptionIDs   = errors.New("invalid option ids")
	ErrExamNotUpdated     = errors.New("exam not updated")
	ErrExamAlreadyStarted = errors.New("can not start two exams at the same time")
)

const (
	defaultTokenTTLSeconds = 24 * 60 * 60
	bearerTokenType        = "Bearer"
	jwtAlgHS256            = "HS256"
	jwtTyp                 = "JWT"
)

type AuthToken struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}

type Student struct {
	StudentID int
	Login     string
	Email     string
	FirstName string
	LastName  string
}

type studentAuth struct {
	StudentID    int
	Login        string
	PasswordHash string
}

type AuthConfig struct {
	JWTSecret     string
	JWTTTLSeconds int
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Reason)
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}

type tokenHeader struct {
	Alg string
	Typ string
}

type tokenClaims struct {
	Sub   string
	Login string
	Exp   int64
	Iat   int64
}

type CourseSummary struct {
	CourseID      int
	Title         string
	TimeLimit     *int
	AvailableType string
	AvailableFrom *string
	AvailableTo   *string
}

type Course struct {
	CourseID      int
	Title         string
	Description   string
	TimeLimit     *int
	AvailableType string
	AvailableFrom *string
	AvailableTo   *string
}

type ExamQuestionRequest struct {
	ExamID     int
	QuestionID int
	StudentID  int
}

type ExamSaveAnswerRequest struct {
	ExamID     int
	QuestionID int
	OptionIDs  []int
	StudentID  int
}

type ExamSubmitRequest struct {
	ExamID    int
	StudentID int
}

type ExamMyListRequest struct {
	StudentID int
	Page      int
	PageSize  int
}

type ExamStart struct {
	ExamID      int
	QuestionIDs []int
	StartedAt   string
	FinishedAt  *string
}

type Question struct {
	QuestionID   int
	QuestionText string
	QuestionType string
	PhotoURL     *string
	Options      []QuestionOption
}

type QuestionOption struct {
	OptionID   int
	OptionText string
	IsCorrect  bool
}

type ExamResult struct {
	ExamID         int
	Status         string
	FinalScore     int
	CorrectAnswers int
	TotalQuestions int
}

type ExamSummary struct {
	ExamID     int
	CourseID   int
	Status     string
	FinalScore int
	FinishedAt string
}

type ExamState struct {
	ExamID      int
	CourseID    int
	Status      string
	QuestionIDs []int
	Answers     []ExamAnswer
}

type ExamAnswer struct {
	QuestionID int
	OptionIDs  []int
}

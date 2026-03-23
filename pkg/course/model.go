package course

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
	ErrCourseNotFound     = errors.New("course not found")
	ErrNoQuestions        = errors.New("course has no questions")
)

const (
	defaultTokenTTLSeconds = 24 * 60 * 60
	bearerTokenType        = "Bearer"
	jwtAlgHS256            = "HS256"
	jwtTyp                 = "JWT"

	ExamStatusInProgress = "in_progress"
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
	CourseId      int
	Title         string
	TimeLimit     *int
	AvailableType string
	AvailableFrom *string
	AvailableTo   *string
}

type Course struct {
	CourseId      int
	Title         string
	Description   string
	TimeLimit     *int
	AvailableType string
	AvailableFrom *string
	AvailableTo   *string
}

type ExamQuestionInput struct {
	ExamID     int
	QuestionID int
	StudentID  int
}

type SaveAnswerInput struct {
	ExamID     int
	QuestionID int
	OptionIDs  []int
	StudentID  int
}

type SubmitExamInput struct {
	ExamID    int
	StudentID int
}

type MyExamListInput struct {
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

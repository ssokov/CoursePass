package coursepass

import (
	"courses/pkg/db"
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

type Student db.Student
type Course db.Course
type Exam db.Exam
type Question db.Question
type ExamAnswer = db.ExamAnswer
type QuestionOption = db.QuestionOption
type VfsFile = db.VfsFile

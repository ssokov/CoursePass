package coursepass

import (
	"errors"
	"fmt"

	"courses/pkg/db"
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

type Student struct {
	db.Student
}

func NewStudent(in *db.Student) *Student {
	if in == nil {
		return nil
	}

	return &Student{
		Student: *in,
	}
}

type Course struct {
	db.Course
}

func NewCourse(in *db.Course) *Course {
	if in == nil {
		return nil
	}

	return &Course{
		Course: *in,
	}
}

type Exam struct {
	db.Exam
}

func NewExam(in *db.Exam) *Exam {
	if in == nil {
		return nil
	}

	return &Exam{
		Exam: *in,
	}
}

type Question struct {
	db.Question
}

func NewQuestion(in *db.Question) *Question {
	if in == nil {
		return nil
	}

	return &Question{
		Question: *in,
	}
}

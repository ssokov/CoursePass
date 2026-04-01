package coursepass

import (
	"errors"

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
	ErrAnswerUnavailable  = errors.New("answer unavailable")
	ErrQuestionNotFound   = errors.New("question not found")
	ErrQuestionNotInExam  = errors.New("question does not belong to exam")
	ErrNoQuestions        = errors.New("coursepass has no questions")
	ErrAnswerAlreadySaved = errors.New("answer already saved")
	ErrInvalidOptionIDs   = errors.New("invalid option ids")
	ErrExamNotUpdated     = errors.New("exam not updated")
	ErrExamAlreadyStarted = errors.New("can not start two exams at the same time")
)

const (
	DefaultTokenTTLSeconds = 24 * 60 * 60
	BearerTokenType        = "Bearer"
	JwtAlgHS256            = "HS256"
	JwtTyp                 = "JWT"
)

type AuthToken struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
}

type TokenHeader struct {
	Alg string
	Typ string
}

type TokenClaims struct {
	Sub   string
	Login string
	Exp   int64
	Iat   int64
}

type StudentDraft struct {
	Login     string `json:"login" validate:"required,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=255"`
	FirstName string `json:"firstName" validate:"required,max=255"`
	LastName  string `json:"lastName" validate:"required,max=255"`
}

type StudentLogin struct {
	Login    string `json:"login" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
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

	Options   []QuestionOption
	PhotoFile *VfsFile
}

func NewQuestion(in *db.Question) *Question {
	if in == nil {
		return nil
	}

	return &Question{
		Question:  *in,
		Options:   NewQuestionOptions(in.Options),
		PhotoFile: NewVfsFile(in.PhotoFile),
	}
}

type QuestionOption struct {
	db.QuestionOption
}

func NewQuestionOption(in *db.QuestionOption) *QuestionOption {
	if in == nil {
		return nil
	}

	return &QuestionOption{
		QuestionOption: *in,
	}
}

type VfsFile struct {
	db.VfsFile
}

func NewVfsFile(in *db.VfsFile) *VfsFile {
	if in == nil {
		return nil
	}

	return &VfsFile{
		VfsFile: *in,
	}
}

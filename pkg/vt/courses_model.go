//nolint:dupl
package vt

import (
	"time"

	"courses/pkg/db"
)

type Course struct {
	ID               int        `json:"id"`
	Title            string     `json:"title" validate:"required,max=255"`
	Description      string     `json:"description" validate:"required"`
	AvailabilityType string     `json:"availabilityType" validate:"required,max=255"`
	AvailableFrom    *time.Time `json:"availableFrom"`
	AvailableTo      *time.Time `json:"availableTo"`
	TimeLimitMinutes *int       `json:"timeLimitMinutes"`
	CreatedAt        time.Time  `json:"createdAt"`
	StatusID         int        `json:"statusId" validate:"required,status"`

	Status *Status `json:"status"`
}

func (c *Course) ToDB() *db.Course {
	if c == nil {
		return nil
	}

	course := &db.Course{
		ID:               c.ID,
		Title:            c.Title,
		Description:      c.Description,
		AvailabilityType: c.AvailabilityType,
		AvailableFrom:    c.AvailableFrom,
		AvailableTo:      c.AvailableTo,
		TimeLimitMinutes: c.TimeLimitMinutes,
		CreatedAt:        c.CreatedAt,
		StatusID:         c.StatusID,
	}

	return course
}

type CourseSearch struct {
	ID               *int       `json:"id"`
	Title            *string    `json:"title"`
	Description      *string    `json:"description"`
	AvailabilityType *string    `json:"availabilityType"`
	AvailableFrom    *time.Time `json:"availableFrom"`
	AvailableTo      *time.Time `json:"availableTo"`
	TimeLimitMinutes *int       `json:"timeLimitMinutes"`
	CreatedAt        *time.Time `json:"createdAt"`
	StatusID         *int       `json:"statusId"`
	IDs              []int      `json:"ids"`
	AvailableToFrom  *time.Time `json:"availableToFrom"`
	AvailableFromTo  *time.Time `json:"availableFromTo"`
}

func (cs *CourseSearch) ToDB() *db.CourseSearch {
	if cs == nil {
		return nil
	}

	return &db.CourseSearch{
		ID:                    cs.ID,
		TitleILike:            cs.Title,
		DescriptionILike:      cs.Description,
		AvailabilityTypeILike: cs.AvailabilityType,
		AvailableFrom:         cs.AvailableFrom,
		AvailableTo:           cs.AvailableTo,
		TimeLimitMinutes:      cs.TimeLimitMinutes,
		CreatedAt:             cs.CreatedAt,
		StatusID:              cs.StatusID,
		IDs:                   cs.IDs,
		AvailableToFrom:       cs.AvailableToFrom,
		AvailableFromTo:       cs.AvailableFromTo,
	}
}

type CourseSummary struct {
	ID               int        `json:"id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	AvailabilityType string     `json:"availabilityType"`
	AvailableFrom    *time.Time `json:"availableFrom"`
	AvailableTo      *time.Time `json:"availableTo"`
	TimeLimitMinutes *int       `json:"timeLimitMinutes"`
	CreatedAt        time.Time  `json:"createdAt"`

	Status *Status `json:"status"`
}

type Exam struct {
	ID             int         `json:"id"`
	CourseID       int         `json:"courseId" validate:"required"`
	StudentID      int         `json:"studentId" validate:"required"`
	Answers        ExamAnswers `json:"answers" validate:"required"`
	TotalQuestions *int        `json:"totalQuestions"`
	CorrectAnswers *int        `json:"correctAnswers"`
	Status         string      `json:"status" validate:"required,max=255"`
	FinalScore     *float64    `json:"finalScore"`
	FinishedAt     *time.Time  `json:"finishedAt"`
	CreatedAt      time.Time   `json:"createdAt"`
	QuestionIDs    []int       `json:"questionIds" validate:"required"`

	Course  *CourseSummary  `json:"course"`
	Student *StudentSummary `json:"student"`
}

func (e *Exam) ToDB() *db.Exam {
	if e == nil {
		return nil
	}

	exam := &db.Exam{
		ID:             e.ID,
		CourseID:       e.CourseID,
		StudentID:      e.StudentID,
		TotalQuestions: e.TotalQuestions,
		CorrectAnswers: e.CorrectAnswers,
		Status:         e.Status,
		FinalScore:     e.FinalScore,
		FinishedAt:     e.FinishedAt,
		CreatedAt:      e.CreatedAt,
		QuestionIDs:    e.QuestionIDs,
	}

	if examAnswers := e.Answers.ToDB(); examAnswers != nil {
		exam.Answers = *examAnswers
	}

	return exam
}

type ExamSearch struct {
	ID             *int       `json:"id"`
	CourseID       *int       `json:"courseId"`
	StudentID      *int       `json:"studentId"`
	TotalQuestions *int       `json:"totalQuestions"`
	CorrectAnswers *int       `json:"correctAnswers"`
	Status         *string    `json:"status"`
	FinalScore     *float64   `json:"finalScore"`
	FinishedAt     *time.Time `json:"finishedAt"`
	CreatedAt      *time.Time `json:"createdAt"`
	IDs            []int      `json:"ids"`
	StatusIn       []string   `json:"statusIn"`
}

func (es *ExamSearch) ToDB() *db.ExamSearch {
	if es == nil {
		return nil
	}

	return &db.ExamSearch{
		ID:             es.ID,
		CourseID:       es.CourseID,
		StudentID:      es.StudentID,
		TotalQuestions: es.TotalQuestions,
		CorrectAnswers: es.CorrectAnswers,
		StatusILike:    es.Status,
		FinalScore:     es.FinalScore,
		FinishedAt:     es.FinishedAt,
		CreatedAt:      es.CreatedAt,
		IDs:            es.IDs,
		StatusIn:       es.StatusIn,
	}
}

type ExamSummary struct {
	ID             int        `json:"id"`
	CourseID       int        `json:"courseId"`
	StudentID      int        `json:"studentId"`
	TotalQuestions *int       `json:"totalQuestions"`
	CorrectAnswers *int       `json:"correctAnswers"`
	Status         string     `json:"status"`
	FinalScore     *float64   `json:"finalScore"`
	FinishedAt     *time.Time `json:"finishedAt"`
	CreatedAt      time.Time  `json:"createdAt"`

	Course  *CourseSummary  `json:"course"`
	Student *StudentSummary `json:"student"`
}

type ExamAnswers struct {
}

func (ea *ExamAnswers) ToDB() *db.ExamAnswers {
	return &db.ExamAnswers{}
}

type Question struct {
	ID           int             `json:"id"`
	CourseID     int             `json:"courseId" validate:"required"`
	PhotoFileID  *int            `json:"photoFileId"`
	QuestionText string          `json:"questionText" validate:"required"`
	Options      QuestionOptions `json:"options" validate:"required"`
	QuestionType string          `json:"questionType" validate:"required,max=255"`
	CreatedAt    time.Time       `json:"createdAt"`

	Course    *CourseSummary  `json:"course"`
	PhotoFile *VfsFileSummary `json:"photoFile"`
}

func (q *Question) ToDB() *db.Question {
	if q == nil {
		return nil
	}

	question := &db.Question{
		ID:           q.ID,
		CourseID:     q.CourseID,
		PhotoFileID:  q.PhotoFileID,
		QuestionText: q.QuestionText,
		QuestionType: q.QuestionType,
		CreatedAt:    q.CreatedAt,
	}

	if questionOptions := q.Options.ToDB(); questionOptions != nil {
		question.Options = *questionOptions
	}

	return question
}

type QuestionSearch struct {
	ID           *int       `json:"id"`
	CourseID     *int       `json:"courseId"`
	PhotoFileID  *int       `json:"photoFileId"`
	QuestionText *string    `json:"questionText"`
	QuestionType *string    `json:"questionType"`
	CreatedAt    *time.Time `json:"createdAt"`
	IDs          []int      `json:"ids"`
}

func (qs *QuestionSearch) ToDB() *db.QuestionSearch {
	if qs == nil {
		return nil
	}

	return &db.QuestionSearch{
		ID:                qs.ID,
		CourseID:          qs.CourseID,
		PhotoFileID:       qs.PhotoFileID,
		QuestionTextILike: qs.QuestionText,
		QuestionTypeILike: qs.QuestionType,
		CreatedAt:         qs.CreatedAt,
		IDs:               qs.IDs,
	}
}

type QuestionSummary struct {
	ID           int       `json:"id"`
	CourseID     int       `json:"courseId"`
	PhotoFileID  *int      `json:"photoFileId"`
	QuestionText string    `json:"questionText"`
	QuestionType string    `json:"questionType"`
	CreatedAt    time.Time `json:"createdAt"`

	Course    *CourseSummary  `json:"course"`
	PhotoFile *VfsFileSummary `json:"photoFile"`
}

type QuestionOptions struct {
}

func (qo *QuestionOptions) ToDB() *db.QuestionOptions {
	return &db.QuestionOptions{}
}

type Student struct {
	ID           int       `json:"id"`
	Login        string    `json:"login" validate:"required,max=255"`
	PasswordHash string    `json:"passwordHash" validate:"required,max=255"`
	FirstName    string    `json:"firstName" validate:"required,max=255"`
	LastName     string    `json:"lastName" validate:"required,max=255"`
	Email        string    `json:"email" validate:"required,email,max=255"`
	CreatedAt    time.Time `json:"createdAt"`
	StatusID     int       `json:"statusId" validate:"required,status"`

	Status *Status `json:"status"`
}

func (s *Student) ToDB() *db.Student {
	if s == nil {
		return nil
	}

	student := &db.Student{
		ID:           s.ID,
		Login:        s.Login,
		PasswordHash: s.PasswordHash,
		FirstName:    s.FirstName,
		LastName:     s.LastName,
		Email:        s.Email,
		CreatedAt:    s.CreatedAt,
		StatusID:     s.StatusID,
	}

	return student
}

type StudentSearch struct {
	ID           *int       `json:"id"`
	Login        *string    `json:"login"`
	PasswordHash *string    `json:"passwordHash"`
	FirstName    *string    `json:"firstName"`
	LastName     *string    `json:"lastName"`
	Email        *string    `json:"email"`
	CreatedAt    *time.Time `json:"createdAt"`
	StatusID     *int       `json:"statusId"`
	IDs          []int      `json:"ids"`
}

func (ss *StudentSearch) ToDB() *db.StudentSearch {
	if ss == nil {
		return nil
	}

	return &db.StudentSearch{
		ID:                ss.ID,
		LoginILike:        ss.Login,
		PasswordHashILike: ss.PasswordHash,
		FirstNameILike:    ss.FirstName,
		LastNameILike:     ss.LastName,
		EmailILike:        ss.Email,
		CreatedAt:         ss.CreatedAt,
		StatusID:          ss.StatusID,
		IDs:               ss.IDs,
	}
}

type StudentSummary struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"passwordHash"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt"`

	Status *Status `json:"status"`
}

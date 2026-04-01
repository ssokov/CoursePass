package rpc

import (
	"path"
	"strings"
	"time"

	"courses/pkg/coursepass"
)

const dateTimeLayout = "2006-01-02 15:04:05"

func formatTimePtr(v *time.Time) *string {
	if v == nil {
		return nil
	}
	s := v.Format(dateTimeLayout)
	return &s
}

func (in StudentDraft) ToModel() coursepass.StudentDraft {
	return coursepass.StudentDraft{
		Login:     in.Login,
		Email:     in.Email,
		Password:  in.Password,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
}

func (in StudentLogin) ToModel() coursepass.StudentLogin {
	return coursepass.StudentLogin{
		Login:    in.Login,
		Password: in.Password,
	}
}

func newFieldError(fe coursepass.FieldError) FieldError {
	var c *FieldErrorConstraint
	if fe.Constraint != nil {
		c = &FieldErrorConstraint{Max: fe.Constraint.Max, Min: fe.Constraint.Min}
	}
	return FieldError{Field: fe.Field, Error: fe.Error, Constraint: c}
}

func newToken(token *coursepass.AuthToken) *Token {
	return &Token{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newStudent(student *coursepass.Student) *Student {
	return &Student{
		StudentID: student.ID,
		Login:     student.Login,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
	}
}

func newCourse(course *coursepass.Course) *Course {
	return &Course{
		CourseID:      course.ID,
		Title:         course.Title,
		Description:   course.Description,
		TimeLimit:     course.TimeLimitMinutes,
		AvailableType: course.AvailabilityType,
		AvailableFrom: formatTimePtr(course.AvailableFrom),
		AvailableTo:   formatTimePtr(course.AvailableTo),
	}
}

func newCourseSummary(course coursepass.Course) CourseSummary {
	return CourseSummary{
		CourseID:      course.ID,
		Title:         course.Title,
		TimeLimit:     course.TimeLimitMinutes,
		AvailableType: course.AvailabilityType,
		AvailableFrom: formatTimePtr(course.AvailableFrom),
		AvailableTo:   formatTimePtr(course.AvailableTo),
	}
}

func newExamStart(exam *coursepass.Exam) *ExamStart {
	return &ExamStart{
		ExamID:      exam.ID,
		QuestionIDs: exam.QuestionIDs,
		StartedAt:   exam.CreatedAt.Format(dateTimeLayout),
		FinishedAt:  formatTimePtr(exam.FinishedAt),
	}
}

func newQuestion(question *coursepass.Question, mediaWebPath string) *Question {
	return &Question{
		QuestionID:   question.ID,
		QuestionText: question.QuestionText,
		QuestionType: question.QuestionType,
		PhotoURL:     newQuestionPhotoURL(question.PhotoFile, mediaWebPath),
		Options:      NewQuestionOptions(question.Options),
	}
}

func newQuestionPhotoURL(photoFile *coursepass.VfsFile, mediaWebPath string) *string {
	if photoFile == nil || photoFile.Path == "" {
		return nil
	}

	basePath := strings.TrimSpace(mediaWebPath)
	if basePath == "" {
		url := photoFile.Path
		return &url
	}

	url := path.Join(basePath, strings.TrimPrefix(photoFile.Path, "/"))
	return &url
}

func NewQuestionOption(option coursepass.QuestionOption) QuestionOption {
	return QuestionOption{
		OptionID:   option.OptionID,
		OptionText: option.OptionText,
	}
}

func newExamResult(exam *coursepass.Exam) *ExamResult {
	finalScore := 0
	if exam.FinalScore != nil {
		finalScore = int(*exam.FinalScore)
	}

	correctAnswers := 0
	if exam.CorrectAnswers != nil {
		correctAnswers = *exam.CorrectAnswers
	}

	totalQuestions := 0
	if exam.TotalQuestions != nil {
		totalQuestions = *exam.TotalQuestions
	}

	return &ExamResult{
		ExamID:         exam.ID,
		Status:         exam.Status,
		FinalScore:     finalScore,
		CorrectAnswers: correctAnswers,
		TotalQuestions: totalQuestions,
	}
}

func newExamSummary(exam coursepass.Exam) ExamSummary {
	finalScore := 0
	if exam.FinalScore != nil {
		finalScore = int(*exam.FinalScore)
	}

	finishedAt := ""
	if exam.FinishedAt != nil {
		finishedAt = exam.FinishedAt.Format(dateTimeLayout)
	}

	return ExamSummary{
		ExamID:     exam.ID,
		CourseID:   exam.CourseID,
		Status:     exam.Status,
		FinalScore: finalScore,
		FinishedAt: finishedAt,
	}
}

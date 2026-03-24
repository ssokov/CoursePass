package vt

import (
	"courses/pkg/db"
)

func NewCourse(in *db.Course) *Course {
	if in == nil {
		return nil
	}

	course := &Course{
		ID:               in.ID,
		Title:            in.Title,
		Description:      in.Description,
		AvailabilityType: in.AvailabilityType,
		AvailableFrom:    in.AvailableFrom,
		AvailableTo:      in.AvailableTo,
		TimeLimitMinutes: in.TimeLimitMinutes,
		CreatedAt:        in.CreatedAt,
		StatusID:         in.StatusID,

		Status: NewStatus(in.StatusID),
	}

	return course
}

func NewCourseSummary(in *db.Course) *CourseSummary {
	if in == nil {
		return nil
	}

	return &CourseSummary{
		ID:               in.ID,
		Title:            in.Title,
		Description:      in.Description,
		AvailabilityType: in.AvailabilityType,
		AvailableFrom:    in.AvailableFrom,
		AvailableTo:      in.AvailableTo,
		TimeLimitMinutes: in.TimeLimitMinutes,
		CreatedAt:        in.CreatedAt,

		Status: NewStatus(in.StatusID),
	}
}

func NewExam(in *db.Exam) *Exam {
	if in == nil {
		return nil
	}

	exam := &Exam{
		ID:             in.ID,
		CourseID:       in.CourseID,
		StudentID:      in.StudentID,
		TotalQuestions: in.TotalQuestions,
		CorrectAnswers: in.CorrectAnswers,
		Status:         in.Status,
		FinalScore:     in.FinalScore,
		FinishedAt:     in.FinishedAt,
		CreatedAt:      in.CreatedAt,
		QuestionIDs:    in.QuestionIDs,

		Course:  NewCourseSummary(in.Course),
		Student: NewStudentSummary(in.Student),
	}

	if examAnswers := NewExamAnswers(&in.Answers); examAnswers != nil {
		exam.Answers = *examAnswers
	}

	return exam
}

func NewExamSummary(in *db.Exam) *ExamSummary {
	if in == nil {
		return nil
	}

	return &ExamSummary{
		ID:             in.ID,
		CourseID:       in.CourseID,
		StudentID:      in.StudentID,
		TotalQuestions: in.TotalQuestions,
		CorrectAnswers: in.CorrectAnswers,
		Status:         in.Status,
		FinalScore:     in.FinalScore,
		FinishedAt:     in.FinishedAt,
		CreatedAt:      in.CreatedAt,

		Course:  NewCourseSummary(in.Course),
		Student: NewStudentSummary(in.Student),
	}
}

func NewExamAnswers(in *db.ExamAnswers) *ExamAnswers {
	return &ExamAnswers{}
}

func NewQuestion(in *db.Question) *Question {
	if in == nil {
		return nil
	}

	question := &Question{
		ID:           in.ID,
		CourseID:     in.CourseID,
		PhotoFileID:  in.PhotoFileID,
		QuestionText: in.QuestionText,
		QuestionType: in.QuestionType,
		CreatedAt:    in.CreatedAt,

		Course:    NewCourseSummary(in.Course),
		PhotoFile: NewVfsFileSummary(in.PhotoFile),
	}

	if questionOptions := NewQuestionOptions(&in.Options); questionOptions != nil {
		question.Options = *questionOptions
	}

	return question
}

func NewQuestionSummary(in *db.Question) *QuestionSummary {
	if in == nil {
		return nil
	}

	return &QuestionSummary{
		ID:           in.ID,
		CourseID:     in.CourseID,
		PhotoFileID:  in.PhotoFileID,
		QuestionText: in.QuestionText,
		QuestionType: in.QuestionType,
		CreatedAt:    in.CreatedAt,

		Course:    NewCourseSummary(in.Course),
		PhotoFile: NewVfsFileSummary(in.PhotoFile),
	}
}

func NewQuestionOptions(in *db.QuestionOptions) *QuestionOptions {
	return &QuestionOptions{}
}

func NewStudent(in *db.Student) *Student {
	if in == nil {
		return nil
	}

	student := &Student{
		ID:           in.ID,
		Login:        in.Login,
		PasswordHash: in.PasswordHash,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Email:        in.Email,
		CreatedAt:    in.CreatedAt,
		StatusID:     in.StatusID,

		Status: NewStatus(in.StatusID),
	}

	return student
}

func NewStudentSummary(in *db.Student) *StudentSummary {
	if in == nil {
		return nil
	}

	return &StudentSummary{
		ID:           in.ID,
		Login:        in.Login,
		PasswordHash: in.PasswordHash,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
		Email:        in.Email,
		CreatedAt:    in.CreatedAt,

		Status: NewStatus(in.StatusID),
	}
}

package rpc

import (
	"courses/pkg/coursepass"
)

func newRegisterResponse(token coursepass.AuthToken) RegisterResponse {
	return RegisterResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newLoginResponse(token coursepass.AuthToken) LoginResponse {
	return LoginResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newMeResponse(student *coursepass.Student) MeResponse {
	return MeResponse{
		StudentID: student.StudentID,
		Login:     student.Login,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
	}
}

func newCourse(c coursepass.Course) Course {
	return Course{
		CourseID:      c.CourseID,
		Title:         c.Title,
		Description:   c.Description,
		TimeLimit:     c.TimeLimit,
		AvailableType: c.AvailableType,
		AvailableFrom: c.AvailableFrom,
		AvailableTo:   c.AvailableTo,
	}
}

func newCourseByIDResponse(course coursepass.Course) ByIDResponse {
	return ByIDResponse{
		Course: newCourse(course),
	}
}

func newCourseSummary(course coursepass.CourseSummary) CourseSummary {
	return CourseSummary{
		CourseID:      course.CourseID,
		Title:         course.Title,
		TimeLimit:     course.TimeLimit,
		AvailableType: course.AvailableType,
		AvailableFrom: course.AvailableFrom,
		AvailableTo:   course.AvailableTo,
	}
}

func newCoursesSummaryResponse(courses []coursepass.CourseSummary) ListResponse {
	result := newCourseSummaries(courses)
	return ListResponse{
		Courses: result,
	}
}

func newExamStartResponse(start coursepass.ExamStart) ExamStartResponse {
	return ExamStartResponse{
		ExamID:      start.ExamID,
		QuestionIDs: start.QuestionIDs,
		StartedAt:   start.StartedAt,
		FinishedAt:  start.FinishedAt,
	}
}

func newQuestionResponse(question coursepass.Question) Question {
	return Question{
		QuestionID:   question.QuestionID,
		QuestionText: question.QuestionText,
		QuestionType: question.QuestionType,
		PhotoURL:     question.PhotoURL,
		Options:      NewQuestionOptions(question.Options),
	}
}

func NewQuestionOption(option coursepass.QuestionOption) QuestionOption {
	return QuestionOption{
		OptionID:   option.OptionID,
		OptionText: option.OptionText,
	}
}

func newExamResultResponse(result coursepass.ExamResult) ExamResult {
	return ExamResult{
		ExamID:         result.ExamID,
		Status:         result.Status,
		FinalScore:     result.FinalScore,
		CorrectAnswers: result.CorrectAnswers,
		TotalQuestions: result.TotalQuestions,
	}
}

func newExamSummary(summary coursepass.ExamSummary) ExamSummary {
	return ExamSummary{
		ExamID:     summary.ExamID,
		CourseID:   summary.CourseID,
		Status:     summary.Status,
		FinalScore: summary.FinalScore,
		FinishedAt: summary.FinishedAt,
	}
}

func newExamHistoryResponse(exams []coursepass.ExamSummary) ExamHistoryResponse {
	return ExamHistoryResponse{
		Exams: newExamSummaries(exams),
	}
}

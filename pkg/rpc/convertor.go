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
	result := make([]CourseSummary, len(courses))
	for i, course := range courses {
		result[i] = newCourseSummary(course)
	}
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

package rpc

import (
	"courses/pkg/course"
)

func newRegisterResponse(token course.AuthToken) RegisterResponse {
	return RegisterResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newLoginResponse(token course.AuthToken) LoginResponse {
	return LoginResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}
}

func newMeResponse(student *course.Student) MeResponse {
	return MeResponse{
		StudentID: student.StudentID,
		Login:     student.Login,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
	}
}

func newCourse(c course.Course) Course {
	return Course{
		CourseId:      c.CourseId,
		Title:         c.Title,
		Description:   c.Description,
		TimeLimit:     c.TimeLimit,
		AvailableType: c.AvailableType,
		AvailableFrom: c.AvailableFrom,
		AvailableTo:   c.AvailableTo,
	}
}

func newCourseByIdResponse(course course.Course) ByIdResponse {
	return ByIdResponse{
		Course: newCourse(course),
	}
}

func newCourseSummary(course course.CourseSummary) CourseSummary {
	return CourseSummary{
		CourseId:      course.CourseId,
		Title:         course.Title,
		TimeLimit:     course.TimeLimit,
		AvailableType: course.AvailableType,
		AvailableFrom: course.AvailableFrom,
		AvailableTo:   course.AvailableTo,
	}
}

func newCoursesSummaryResponse(courses []course.CourseSummary) ListResponse {
	result := make([]CourseSummary, len(courses))
	for i, course := range courses {
		result[i] = newCourseSummary(course)
	}
	return ListResponse{
		Courses: result,
	}
}

func newExamStartResponse(start course.ExamStart) ExamStartResponse {
	return ExamStartResponse{
		ExamID:      start.ExamID,
		QuestionIDs: start.QuestionIDs,
		StartedAt:   start.StartedAt,
		FinishedAt:  start.FinishedAt,
	}
}

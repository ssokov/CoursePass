package coursepass

import (
	"path"
	"strconv"
	"strings"
	"time"

	"courses/pkg/db"
)

const dateTimeLayout = "2006-01-02 15:04:05"

func newDBStudent(login, passwordHash, firstName, lastName, email string) *db.Student {
	return &db.Student{
		Login:        login,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		StatusID:     db.StatusEnabled,
	}
}

func newStudent(student db.Student) Student {
	return Student{
		StudentID: student.ID,
		Login:     student.Login,
		Email:     student.Email,
		FirstName: student.FirstName,
		LastName:  student.LastName,
	}
}

func newAuthToken(token string, expiresIn int) AuthToken {
	return AuthToken{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   bearerTokenType,
	}
}

func newTokenHeader() tokenHeader {
	return tokenHeader{
		Alg: jwtAlgHS256,
		Typ: jwtTyp,
	}
}

func newTokenClaims(studentID int, login string, iat, exp int64) tokenClaims {
	return tokenClaims{
		Sub:   strconv.Itoa(studentID),
		Login: login,
		Exp:   exp,
		Iat:   iat,
	}
}

func newCourseSummary(course db.Course) CourseSummary {
	return CourseSummary{
		CourseID:      course.ID,
		Title:         course.Title,
		TimeLimit:     course.TimeLimitMinutes,
		AvailableType: course.AvailabilityType,
		AvailableFrom: formatTimePtr(course.AvailableFrom),
		AvailableTo:   formatTimePtr(course.AvailableTo),
	}
}

func newCourse(course db.Course) Course {
	return Course{
		CourseID:      course.ID,
		Title:         course.Title,
		Description:   course.Description,
		TimeLimit:     course.TimeLimitMinutes,
		AvailableType: course.AvailabilityType,
		AvailableFrom: formatTimePtr(course.AvailableFrom),
		AvailableTo:   formatTimePtr(course.AvailableTo),
	}
}

func formatTimePtr(v *time.Time) *string {
	if v == nil {
		return nil
	}
	s := v.Format(dateTimeLayout)
	return &s
}

func newExamStart(exam db.Exam, questionIDs []int) ExamStart {
	return ExamStart{
		ExamID:      exam.ID,
		QuestionIDs: questionIDs,
		StartedAt:   exam.CreatedAt.Format(dateTimeLayout),
		FinishedAt:  formatTimePtr(exam.FinishedAt),
	}
}

func newQuestion(question db.Question, mediaWebPath string) Question {
	return Question{
		QuestionID:   question.ID,
		QuestionText: question.QuestionText,
		QuestionType: question.QuestionType,
		PhotoURL:     newQuestionPhotoURL(question.PhotoFile, mediaWebPath),
		Options:      newQuestionOptions(question.Options),
	}
}

func newQuestionOptions(options db.QuestionOptions) []QuestionOption {
	result := make([]QuestionOption, len(options))
	for i := range options {
		result[i] = QuestionOption{
			OptionID:   options[i].OptionID,
			OptionText: options[i].OptionText,
		}
	}

	return result
}

func newQuestionPhotoURL(photoFile *db.VfsFile, mediaWebPath string) *string {
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

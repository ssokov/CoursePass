package coursepass

import (
	"strconv"
	"time"

	"courses/pkg/db"
)

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

func newStudent(student *db.Student) *Student {
	if student == nil {
		return nil
	}
	s := Student(*student)
	return &s
}

func newAuthToken(token string, expiresIn int) *AuthToken {
	return &AuthToken{
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

func newCourse(course *db.Course) *Course {
	if course == nil {
		return nil
	}
	c := Course(*course)
	return &c
}

func NewCourse(course db.Course) Course {
	return Course(course)
}

func newExam(e *db.Exam) *Exam {
	if e == nil {
		return nil
	}
	r := Exam(*e)
	return &r
}

func NewExam(exam db.Exam) Exam {
	return Exam(exam)
}

func newQuestion(q *db.Question) *Question {
	if q == nil {
		return nil
	}
	r := Question(*q)
	return &r
}

func NewQuestion(question db.Question) Question {
	return Question(question)
}

func newDBExamAnswersUpdate(examID int, answers db.ExamAnswers) *db.Exam {
	return &db.Exam{
		ID:      examID,
		Answers: answers,
	}
}

func newDBExamSubmitUpdate(examID int, status string, correctAnswers, totalQuestions int, finalScore float64, finishedAt time.Time) *db.Exam {
	return &db.Exam{
		ID:             examID,
		Status:         status,
		CorrectAnswers: &correctAnswers,
		TotalQuestions: &totalQuestions,
		FinalScore:     &finalScore,
		FinishedAt:     &finishedAt,
	}
}

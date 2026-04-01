package coursepass

import (
	"strconv"
	"time"

	"courses/pkg/db"
)

func NewDBStudent(login, passwordHash, firstName, lastName, email string) *db.Student {
	return &db.Student{
		Login:        login,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		StatusID:     db.StatusEnabled,
	}
}

func NewAuthToken(token string, expiresIn int) *AuthToken {
	return &AuthToken{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   BearerTokenType,
	}
}

func NewTokenHeader() TokenHeader {
	return TokenHeader{
		Alg: JwtAlgHS256,
		Typ: JwtTyp,
	}
}

func NewTokenClaims(studentID int, login string, iat, exp int64) TokenClaims {
	return TokenClaims{
		Sub:   strconv.Itoa(studentID),
		Login: login,
		Exp:   exp,
		Iat:   iat,
	}
}

func NewDBExamAnswersUpdate(examID int, answers db.ExamAnswers) *db.Exam {
	return &db.Exam{
		ID:      examID,
		Answers: answers,
	}
}

func NewDBExamSubmitUpdate(examID int, status string, correctAnswers, totalQuestions int, finalScore float64, finishedAt time.Time) *db.Exam {
	return &db.Exam{
		ID:             examID,
		Status:         status,
		CorrectAnswers: &correctAnswers,
		TotalQuestions: &totalQuestions,
		FinalScore:     &finalScore,
		FinishedAt:     &finishedAt,
	}
}

package rpc

import (
	"courses/pkg/coursepass"
	"errors"

	"github.com/vmkteam/zenrpc/v2"
)

const (
	errInvalidCredentials = -32001
	errInvalidToken       = -32002
)

func mapRPCError(err error) error {
	var validationErr coursepass.ValidationError
	switch {
	case errors.As(err, &validationErr):
		return invalidParamsError(validationErr.Field, validationErr.Reason)
	case errors.Is(err, coursepass.ErrLoginExists):
		return invalidParamsError("login", "must be unique")
	case errors.Is(err, coursepass.ErrEmailExists):
		return invalidParamsError("email", "must be unique")
	case errors.Is(err, coursepass.ErrInvalidCredentials):
		return &zenrpc.Error{
			Code:    errInvalidCredentials,
			Message: "invalid credentials",
		}
	case errors.Is(err, coursepass.ErrInvalidToken):
		return &zenrpc.Error{
			Code:    errInvalidToken,
			Message: "invalid token",
		}
	case errors.Is(err, coursepass.ErrStudentNotFound):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "student not found",
		}
	case errors.Is(err, coursepass.ErrCourseNotFound):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "coursepass not found",
		}
	case errors.Is(err, coursepass.ErrExamNotFound):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "exam not found",
		}
	case errors.Is(err, coursepass.ErrExamNotInProgress):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "exam is not in progress",
		}
	case errors.Is(err, coursepass.ErrQuestionNotFound):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "question not found",
		}
	case errors.Is(err, coursepass.ErrNoQuestions):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "coursepass has no questions",
		}
	default:
		return zenrpc.NewError(zenrpc.InternalError, err)
	}
}

func invalidParamsError(field, reason string) *zenrpc.Error {
	return &zenrpc.Error{
		Code:    zenrpc.InvalidParams,
		Message: "Invalid params",
		Data: map[string]any{
			"field":  field,
			"reason": reason,
		},
	}
}

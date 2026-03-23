package rpc

import (
	"courses/pkg/course"
	"errors"

	"github.com/vmkteam/zenrpc/v2"
)

const (
	errInvalidCredentials = -32001
	errInvalidToken       = -32002
)

func mapRPCError(err error) error {
	var validationErr course.ValidationError
	switch {
	case errors.As(err, &validationErr):
		return invalidParamsError(validationErr.Field, validationErr.Reason)
	case errors.Is(err, course.ErrLoginExists):
		return invalidParamsError("login", "must be unique")
	case errors.Is(err, course.ErrEmailExists):
		return invalidParamsError("email", "must be unique")
	case errors.Is(err, course.ErrInvalidCredentials):
		return &zenrpc.Error{
			Code:    errInvalidCredentials,
			Message: "invalid credentials",
		}
	case errors.Is(err, course.ErrInvalidToken):
		return &zenrpc.Error{
			Code:    errInvalidToken,
			Message: "invalid token",
		}
	case errors.Is(err, course.ErrStudentNotFound):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "student not found",
		}
	case errors.Is(err, course.ErrCourseNotFound):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "course not found",
		}
	case errors.Is(err, course.ErrNoQuestions):
		return &zenrpc.Error{
			Code:    zenrpc.InvalidParams,
			Message: "course has no questions",
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

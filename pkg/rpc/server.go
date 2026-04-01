package rpc

import (
	"net/http"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	zm "github.com/vmkteam/zenrpc-middleware"
	"github.com/vmkteam/zenrpc/v2"
)

var (
	ErrNotImplemented = zenrpc.NewStringError(http.StatusInternalServerError, "not implemented")
	ErrInternal       = zenrpc.NewStringError(http.StatusInternalServerError, "internal error")
	ErrNotFound       = zenrpc.NewStringError(http.StatusNotFound, "not found")

	ErrInvalidParams      = zenrpc.NewStringError(zenrpc.InvalidParams, "invalid params")
	ErrInvalidToken       = zenrpc.NewStringError(http.StatusUnauthorized, "invalid token")
	ErrInvalidCredentials = zenrpc.NewStringError(http.StatusUnauthorized, "invalid credentials")
	ErrLoginExists        = zenrpc.NewStringError(http.StatusConflict, "login exists")
	ErrEmailExists        = zenrpc.NewStringError(http.StatusConflict, "email exists")
	ErrExamConflict       = zenrpc.NewStringError(http.StatusConflict, "exam conflict")
)

const (
	NSAuth   = "auth"
	NSCourse = "course"
	NSExam   = "exam"
)

var allowDebugFn = func() zm.AllowDebugFunc {
	return func(req *http.Request) bool {
		return req != nil && req.FormValue("__level") == "5"
	}
}

//go:generate go tool zenrpc

// New returns new zenrpc Server.
func New(dbo db.DB, logger embedlog.Logger, jwtSecret string, jwtTTLSeconds int, isDevel bool, mediaWebPath string) *zenrpc.Server {
	rpc := zenrpc.NewServer(zenrpc.Options{
		ExposeSMD: true,
		AllowCORS: true,
	})

	rpc.Use(
		zm.WithDevel(isDevel),
		zm.WithHeaders(),
		zm.WithSentry(zm.DefaultServerName),
		zm.WithNoCancelContext(),
		zm.WithMetrics(zm.DefaultServerName),
		zm.WithTiming(isDevel, allowDebugFn()),
		zm.WithSQLLogger(dbo.DB, isDevel, allowDebugFn(), allowDebugFn()),
	)

	rpc.Use(
		zm.WithSLog(logger.Print, zm.DefaultServerName, nil),
		zm.WithErrorSLog(logger.Print, zm.DefaultServerName, nil),
		authMiddleware(jwtSecret, logger),
	)

	// services
	rpc.RegisterAll(map[string]zenrpc.Invoker{
		NSAuth:   NewAuthService(dbo, logger, jwtSecret, jwtTTLSeconds),
		NSCourse: NewCoursesService(dbo, logger),
		NSExam:   NewExamService(dbo, logger, mediaWebPath),
	})

	return rpc
}

func newInternalError(err error) *zenrpc.Error {
	return zenrpc.NewError(http.StatusInternalServerError, err)
}

func newInvalidParamsError(field, reason string) *zenrpc.Error {
	return &zenrpc.Error{
		Code:    ErrInvalidParams.Code,
		Message: ErrInvalidParams.Message,
		Data: map[string]any{
			"field":  field,
			"reason": reason,
		},
	}
}

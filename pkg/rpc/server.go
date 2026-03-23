package rpc

import (
	"net/http"

	"courses/pkg/course"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	zm "github.com/vmkteam/zenrpc-middleware"
	"github.com/vmkteam/zenrpc/v2"
)

var (
	ErrNotImplemented = zenrpc.NewStringError(http.StatusInternalServerError, "not implemented")
	ErrInternal       = zenrpc.NewStringError(http.StatusInternalServerError, "internal error")
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
func New(dbo db.DB, logger embedlog.Logger, authCfg course.AuthConfig, isDevel bool) *zenrpc.Server {
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
		authMiddleware(authCfg, logger),
	)

	// services
	rpc.RegisterAll(map[string]zenrpc.Invoker{
		NSAuth:   NewAuthService(dbo, logger, authCfg),
		NSCourse: NewCoursesService(dbo, logger, authCfg),
		NSExam:   NewExamService(dbo, logger, authCfg),
	})

	return rpc
}

//nolint:unused
func newInternalError(err error) *zenrpc.Error {
	return zenrpc.NewError(http.StatusInternalServerError, err)
}

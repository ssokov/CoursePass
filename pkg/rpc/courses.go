package rpc

import (
	"context"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type CoursesService struct {
	zenrpc.Service
	embedlog.Logger

	authManager   *coursepass.AuthManager
	courseManager *coursepass.CourseManager
}

func NewCoursesService(dbc db.DB, logger embedlog.Logger, authCfg coursepass.AuthConfig) *CoursesService {
	return &CoursesService{
		authManager:   coursepass.NewAuthManager(dbc, logger, authCfg),
		courseManager: coursepass.NewCourseManager(dbc, logger),
		Logger:        logger,
	}
}

func (cs *CoursesService) Me(ctx context.Context) (MeResponse, error) {
	studentID, ok := StudentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		cs.Logger.Error(ctx, "coursepass me failed: no studentID in context")
		return MeResponse{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	student, err := cs.authManager.Me(ctx, studentID)
	if err != nil {
		cs.Logger.Error(ctx, "coursepass me failed", "err", err)
		return MeResponse{}, mapRPCError(err)
	}

	return newMeResponse(student), nil
}

func (cs *CoursesService) List(ctx context.Context, req ListRequest) (ListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	courses, err := cs.courseManager.Summary(ctx, req.Page, req.PageSize)
	if err != nil {
		cs.Logger.Error(ctx, "coursepass list failed", "err", err)
		return ListResponse{}, mapRPCError(err)
	}

	return newCoursesSummaryResponse(courses), nil
}

func (cs *CoursesService) ById(ctx context.Context, req ByIDRequest) (ByIDResponse, error) {
	if req.CourseID < 1 {
		return ByIDResponse{}, invalidParamsError("courseId", "must be greater than 0")
	}

	courseObj, err := cs.courseManager.ByID(ctx, req.CourseID)
	if err != nil {
		cs.Logger.Error(ctx, "coursepass by id failed", "err", err)
		return ByIDResponse{}, mapRPCError(err)
	}

	return newCourseByIDResponse(courseObj), nil
}

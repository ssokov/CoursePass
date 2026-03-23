package rpc

import (
	"context"

	"courses/pkg/course"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type ExamService struct {
	zenrpc.Service
	embedlog.Logger

	courseManager *course.CourseManager
}

func NewExamService(dbc db.DB, logger embedlog.Logger, authCfg course.AuthConfig) *ExamService {
	return &ExamService{
		courseManager: course.NewCourseManager(dbc, logger, authCfg),
		Logger:        logger,
	}
}

func (es *ExamService) Start(ctx context.Context, req ExamStartRequest) (ExamStartResponse, error) {
	studentID, ok := StudentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam start failed: no studentID in context")
		return ExamStartResponse{}, mapRPCError(course.ErrInvalidToken)
	}

	start, err := es.courseManager.StartExam(ctx, req.CourseID, studentID)
	if err != nil {
		es.Logger.Error(ctx, "exam start failed", "err", err)
		return ExamStartResponse{}, mapRPCError(err)
	}

	return newExamStartResponse(start), nil
}

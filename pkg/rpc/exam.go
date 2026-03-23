package rpc

import (
	"context"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type ExamService struct {
	zenrpc.Service
	embedlog.Logger

	examManager *coursepass.ExamManager
}

func NewExamService(dbc db.DB, logger embedlog.Logger, authCfg coursepass.AuthConfig) *ExamService {
	_ = authCfg
	return &ExamService{
		examManager: coursepass.NewExamManager(dbc, logger),
		Logger:      logger,
	}
}

func (es *ExamService) Start(ctx context.Context, req ExamStartRequest) (ExamStartResponse, error) {
	studentID, ok := StudentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam start failed: no studentID in context")
		return ExamStartResponse{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	start, err := es.examManager.Start(ctx, req.CourseID, studentID)
	if err != nil {
		es.Logger.Error(ctx, "exam start failed", "err", err)
		return ExamStartResponse{}, mapRPCError(err)
	}

	return newExamStartResponse(start), nil
}

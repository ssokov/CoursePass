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

func NewExamService(dbc db.DB, logger embedlog.Logger, mediaWebPath string) *ExamService {
	return &ExamService{
		examManager: coursepass.NewExamManager(dbc, logger, mediaWebPath),
		Logger:      logger,
	}
}

func (es *ExamService) Start(ctx context.Context, req ExamStartRequest) (ExamStartResponse, error) {
	if req.CourseID < 1 {
		return ExamStartResponse{}, invalidParamsError("courseId", "must be greater than 0")
	}

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

func (es *ExamService) Question(ctx context.Context, req ExamQuestionRequest) (Question, error) {
	if req.ExamID < 1 {
		return Question{}, invalidParamsError("examId", "must be greater than 0")
	}
	if req.QuestionID < 1 {
		return Question{}, invalidParamsError("questionId", "must be greater than 0")
	}

	studentID, ok := StudentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam question failed: no studentID in context")
		return Question{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	question, err := es.examManager.Question(ctx, req.ExamID, req.QuestionID, studentID)
	if err != nil {
		es.Logger.Error(ctx, "exam question failed", "err", err)
		return Question{}, mapRPCError(err)
	}

	return newQuestionResponse(question), nil
}

func (es *ExamService) SaveAnswer(ctx context.Context, req SaveAnswerRequest) error {
	if req.ExamID < 1 {
		return invalidParamsError("examId", "must be greater than 0")
	}
	if req.QuestionID < 1 {
		return invalidParamsError("questionId", "must be greater than 0")
	}
	if len(req.OptionIDs) < 1 {
		return invalidParamsError("optionIds", "size must be bigger than 0")
	}
	studentID, ok := StudentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam save failed: no studentID in context")
		return mapRPCError(coursepass.ErrInvalidToken)
	}

	err := es.examManager.SaveAnswer(ctx, studentID, req.ExamID, req.QuestionID, req.OptionIDs)
	if err != nil {
		es.Logger.Error(ctx, "exam save failed", "err", err)
		return mapRPCError(err)
	}

	return nil
}

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

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam start failed: no studentID in context")
		return ExamStartResponse{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	start, err := es.examManager.Start(ctx, studentID, req.CourseID)
	if err != nil {
		es.Logger.Error(ctx, "exam start failed", "err", err)
		return ExamStartResponse{}, mapRPCError(err)
	}

	return newExamStartResponse(start), nil
}

func (es *ExamService) GetQuestion(ctx context.Context, req ExamGetQuestionRequest) (Question, error) {
	if req.ExamID < 1 {
		return Question{}, invalidParamsError("examId", "must be greater than 0")
	}
	if req.QuestionID < 1 {
		return Question{}, invalidParamsError("questionId", "must be greater than 0")
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam question failed: no studentID in context")
		return Question{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	question, err := es.examManager.Question(ctx, studentID, req.QuestionID, req.ExamID)
	if err != nil {
		es.Logger.Error(ctx, "exam question failed", "err", err)
		return Question{}, mapRPCError(err)
	}

	return newQuestionResponse(question), nil
}

func (es *ExamService) Answer(ctx context.Context, req AnswerRequest) error {
	if req.ExamID < 1 {
		return invalidParamsError("examId", "must be greater than 0")
	}
	if req.QuestionID < 1 {
		return invalidParamsError("questionId", "must be greater than 0")
	}
	if len(req.OptionIDs) < 1 {
		return invalidParamsError("optionIds", "size must be bigger than 0")
	}
	studentID, ok := studentIDFromContext(ctx)
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

func (es *ExamService) Submit(ctx context.Context, req ExamSubmitRequest) (ExamResult, error) {
	if req.ExamID < 1 {
		return ExamResult{}, invalidParamsError("examId", "must be greater than 0")
	}
	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam submit failed: no studentID in context")
		return ExamResult{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	result, err := es.examManager.Submit(ctx, studentID, req.ExamID)
	if err != nil {
		es.Logger.Error(ctx, "exam submit failed", "err", err)
		return ExamResult{}, mapRPCError(err)
	}

	return newExamResultResponse(result), nil
}

func (es *ExamService) History(ctx context.Context, req ExamHistoryRequest) (ExamHistoryResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		es.Logger.Error(ctx, "exam history failed: no studentID in context")
		return ExamHistoryResponse{}, mapRPCError(coursepass.ErrInvalidToken)
	}

	exams, err := es.examManager.MyList(ctx, studentID, req.Page, req.PageSize)
	if err != nil {
		es.Logger.Error(ctx, "exam history failed", "err", err)
		return ExamHistoryResponse{}, mapRPCError(err)
	}

	return newExamHistoryResponse(exams), nil
}

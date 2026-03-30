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

	examManager  *coursepass.ExamManager
	mediaWebPath string
}

func NewExamService(dbc db.DB, logger embedlog.Logger, mediaWebPath string) *ExamService {
	return &ExamService{
		examManager:  coursepass.NewExamManager(dbc, logger, mediaWebPath),
		Logger:       logger,
		mediaWebPath: mediaWebPath,
	}
}

//zenrpc:401 invalid token
func (es *ExamService) Start(ctx context.Context, courseID int) (*ExamStart, error) {
	if courseID < 1 {
		return nil, newInvalidParamsError("courseID", "must be greater than 0")
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return nil, ErrInvalidToken
	}

	exam, err := es.examManager.Start(ctx, studentID, courseID)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newExamStart(exam), nil
}

//zenrpc:401 invalid token
func (es *ExamService) GetQuestion(ctx context.Context, examID, questionID int) (*Question, error) {
	if examID < 1 {
		return nil, newInvalidParamsError("examId", "must be greater than 0")
	}
	if questionID < 1 {
		return nil, newInvalidParamsError("questionId", "must be greater than 0")
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return nil, ErrInvalidToken
	}

	question, err := es.examManager.Question(ctx, studentID, questionID, examID)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newQuestion(question, es.mediaWebPath), nil
}

//zenrpc:401 invalid token
func (es *ExamService) Answer(ctx context.Context, examID, questionID int, optionIDs []int) error {
	if examID < 1 {
		return newInvalidParamsError("examId", "must be greater than 0")
	}
	if questionID < 1 {
		return newInvalidParamsError("questionId", "must be greater than 0")
	}
	if len(optionIDs) < 1 {
		return newInvalidParamsError("optionIds", "size must be bigger than 0")
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return ErrInvalidToken
	}

	err := es.examManager.SaveAnswer(ctx, studentID, examID, questionID, optionIDs)
	if err != nil {
		return newInternalError(err)
	}

	return nil
}

//zenrpc:401 invalid token
func (es *ExamService) Submit(ctx context.Context, examID int) (*ExamResult, error) {
	if examID < 1 {
		return nil, newInvalidParamsError("examId", "must be greater than 0")
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return nil, ErrInvalidToken
	}

	exam, err := es.examManager.Submit(ctx, studentID, examID)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newExamResult(exam), nil
}

//zenrpc:401 invalid token
func (es *ExamService) History(ctx context.Context, page, pageSize int) ([]ExamSummary, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return nil, ErrInvalidToken
	}

	exams, err := es.examManager.MyList(ctx, studentID, page, pageSize)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newExamSummaries(exams), nil
}

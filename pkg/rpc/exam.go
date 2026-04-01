package rpc

import (
	"context"
	"errors"

	"courses/pkg/coursepass"
	"courses/pkg/coursepass/exam"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type ExamService struct {
	zenrpc.Service
	embedlog.Logger

	examManager  *exam.Manager
	mediaWebPath string
}

func NewExamService(dbc db.DB, logger embedlog.Logger, mediaWebPath string) *ExamService {
	return &ExamService{
		examManager:  exam.NewManager(dbc, logger, mediaWebPath),
		Logger:       logger,
		mediaWebPath: mediaWebPath,
	}
}

//zenrpc:401 invalid token
//zenrpc:404 not found
//zenrpc:409 exam conflict
//zenrpc:500 internal error
func (es *ExamService) Start(ctx context.Context, courseID int) (*ExamStart, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidToken
	}

	exam, err := es.examManager.Start(ctx, studentID, courseID)
	if err != nil {
		if errors.Is(err, coursepass.ErrCourseNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, coursepass.ErrNoQuestions) || errors.Is(err, coursepass.ErrExamAlreadyStarted) {
			return nil, ErrExamConflict
		}
		return nil, newInternalError(err)
	}

	return newExamStart(exam), nil
}

//zenrpc:401 invalid token
//zenrpc:404 not found
//zenrpc:409 exam conflict
//zenrpc:500 internal error
func (es *ExamService) GetQuestion(ctx context.Context, examID, questionID int) (*Question, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidToken
	}

	question, err := es.examManager.Question(ctx, studentID, questionID, examID)
	if err != nil {
		if errors.Is(err, coursepass.ErrExamNotFound) || errors.Is(err, coursepass.ErrQuestionNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, coursepass.ErrExamNotInProgress) {
			return nil, ErrExamConflict
		}
		if errors.Is(err, coursepass.ErrQuestionNotInExam) {
			return nil, ErrInvalidParams
		}
		return nil, newInternalError(err)
	}

	return newQuestion(question, es.mediaWebPath), nil
}

//zenrpc:401 invalid token
//zenrpc:404 not found
//zenrpc:409 exam conflict
//zenrpc:500 internal error
func (es *ExamService) Answer(ctx context.Context, examID, questionID int, optionIDs []int) error {
	studentID, ok := studentIDFromContext(ctx)
	if !ok {
		return ErrInvalidToken
	}

	err := es.examManager.SaveAnswer(ctx, studentID, examID, questionID, optionIDs)
	if err != nil {
		if errors.Is(err, coursepass.ErrExamNotFound) {
			return ErrNotFound
		}
		if errors.Is(err, coursepass.ErrAnswerUnavailable) {
			return ErrExamConflict
		}
		if errors.Is(err, coursepass.ErrInvalidOptionIDs) {
			return ErrInvalidParams
		}
		return newInternalError(err)
	}

	return nil
}

//zenrpc:401 invalid token
//zenrpc:404 not found
//zenrpc:409 exam conflict
//zenrpc:500 internal error
func (es *ExamService) Submit(ctx context.Context, examID int) (*ExamResult, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidToken
	}

	exam, err := es.examManager.Submit(ctx, studentID, examID)
	if err != nil {
		if errors.Is(err, coursepass.ErrExamNotFound) {
			return nil, ErrNotFound
		}
		if errors.Is(err, coursepass.ErrExamNotInProgress) || errors.Is(err, coursepass.ErrNoQuestions) {
			return nil, ErrExamConflict
		}
		return nil, newInternalError(err)
	}

	return newExamResult(exam), nil
}

//zenrpc:401 invalid token
//zenrpc:500 internal error
func (es *ExamService) History(ctx context.Context, page, pageSize int) ([]ExamSummary, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidToken
	}

	exams, err := es.examManager.MyList(ctx, studentID, page, pageSize)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newExamSummaries(exams), nil
}

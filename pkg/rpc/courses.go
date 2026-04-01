package rpc

import (
	"context"
	"errors"

	"courses/pkg/coursepass"
	"courses/pkg/coursepass/course"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type CoursesService struct {
	zenrpc.Service
	embedlog.Logger

	courseManager *course.Manager
}

func NewCoursesService(dbc db.DB, logger embedlog.Logger) *CoursesService {
	return &CoursesService{
		courseManager: course.NewManager(dbc, logger),
		Logger:        logger,
	}
}

//zenrpc:401 invalid token
//zenrpc:404 not found
//zenrpc:500 internal error
func (cs *CoursesService) Me(ctx context.Context) (*Student, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok {
		return nil, ErrInvalidToken
	}

	student, err := cs.courseManager.Me(ctx, studentID)
	if err != nil {
		if errors.Is(err, coursepass.ErrStudentNotFound) {
			return nil, ErrNotFound
		}
		return nil, newInternalError(err)
	}
	if student == nil {
		return nil, ErrNotFound
	}

	return newStudent(student), nil
}

//zenrpc:500 internal error
func (cs *CoursesService) List(ctx context.Context, page, pageSize int) ([]CourseSummary, error) {
	courses, err := cs.courseManager.List(ctx, page, pageSize)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newCourseSummaries(courses), nil
}

//zenrpc:404 not found
//zenrpc:500 internal error
func (cs *CoursesService) ByID(ctx context.Context, courseID int) (*Course, error) {
	course, err := cs.courseManager.ByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, coursepass.ErrCourseNotFound) {
			return nil, ErrNotFound
		}
		return nil, newInternalError(err)
	}
	if course == nil {
		return nil, ErrNotFound
	}

	return newCourse(course), nil
}

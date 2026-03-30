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

	courseManager *coursepass.CourseManager
}

func NewCoursesService(dbc db.DB, logger embedlog.Logger) *CoursesService {
	return &CoursesService{
		courseManager: coursepass.NewCourseManager(dbc, logger),
		Logger:        logger,
	}
}

//zenrpc:401 invalid token
//zenrpc:404 not found
func (cs *CoursesService) Me(ctx context.Context) (*Student, error) {
	studentID, ok := studentIDFromContext(ctx)
	if !ok || studentID <= 0 {
		return nil, ErrInvalidToken
	}

	student, err := cs.courseManager.Me(ctx, studentID)
	if student == nil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, newInternalError(err)
	}

	return newStudent(student), nil
}

func (cs *CoursesService) List(ctx context.Context, page, pageSize int) ([]CourseSummary, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	courses, err := cs.courseManager.List(ctx, page, pageSize)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newCourseSummaries(courses), nil
}

//zenrpc:404 not found
func (cs *CoursesService) ByID(ctx context.Context, courseID int) (*Course, error) {
	if courseID < 1 {
		return nil, newInvalidParamsError("courseId", "must be greater than 0")
	}
	course, err := cs.courseManager.ByID(ctx, courseID)
	if course == nil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, newInternalError(err)
	}

	return newCourse(course), nil
}

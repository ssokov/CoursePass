package coursepass

import (
	"context"
	"fmt"
	"time"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
)

type CourseManager struct {
	repo db.CoursesRepo
	embedlog.Logger
}

func NewCourseManager(dbo db.DB, logger embedlog.Logger) *CourseManager {
	return &CourseManager{
		repo:   db.NewCoursesRepo(dbo),
		Logger: logger,
	}
}

func (cm *CourseManager) List(ctx context.Context, page, pageSize int) ([]Course, error) {
	return cm.availableCourses(ctx, time.Now(), page, pageSize)
}

func (cm *CourseManager) ByID(ctx context.Context, courseID int) (*Course, error) {
	return cm.courseByID(ctx, courseID)
}

func (cm *CourseManager) Me(ctx context.Context, studentID int) (*Student, error) {
	return cm.studentByID(ctx, studentID)
}

func (cm *CourseManager) availableCourses(ctx context.Context, currentTime time.Time, page, pageSize int) ([]Course, error) {
	courses, err := cm.repo.CoursesByFilters(ctx, &db.CourseSearch{
		AvailableFromTo: &currentTime,
		AvailableToFrom: &currentTime,
	}, db.Pager{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get courses: %w", err)
	}

	return NewCourses(courses), nil
}

func (cm *CourseManager) courseByID(ctx context.Context, courseID int) (*Course, error) {
	courseData, err := cm.repo.OneCourse(ctx, &db.CourseSearch{
		ID: &courseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get course: %w", err)
	}
	if courseData == nil {
		return nil, ErrCourseNotFound
	}

	return NewCourse(courseData), nil
}

func (cm *CourseManager) studentByID(ctx context.Context, studentID int) (*Student, error) {
	studentData, err := cm.repo.OneStudent(ctx, &db.StudentSearch{
		ID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return nil, ErrStudentNotFound
	}

	return NewStudent(studentData), nil
}

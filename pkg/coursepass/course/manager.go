package course

import (
	"context"
	"fmt"
	"time"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
)

type Manager struct {
	repo db.CoursesRepo
	embedlog.Logger
}

func NewManager(dbo db.DB, logger embedlog.Logger) *Manager {
	return &Manager{
		repo:   db.NewCoursesRepo(dbo),
		Logger: logger,
	}
}

func (cm *Manager) List(ctx context.Context, page, pageSize int) ([]coursepass.Course, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return cm.availableCourses(ctx, time.Now(), page, pageSize)
}

func (cm *Manager) ByID(ctx context.Context, courseID int) (*coursepass.Course, error) {
	return cm.courseByID(ctx, courseID)
}

func (cm *Manager) Me(ctx context.Context, studentID int) (*coursepass.Student, error) {
	return cm.studentByID(ctx, studentID)
}

func (cm *Manager) availableCourses(ctx context.Context, currentTime time.Time, page, pageSize int) ([]coursepass.Course, error) {
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

	return coursepass.NewCourses(courses), nil
}

func (cm *Manager) courseByID(ctx context.Context, courseID int) (*coursepass.Course, error) {
	courseData, err := cm.repo.OneCourse(ctx, &db.CourseSearch{
		ID: &courseID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get course: %w", err)
	}
	if courseData == nil {
		return nil, coursepass.ErrCourseNotFound
	}

	return coursepass.NewCourse(courseData), nil
}

func (cm *Manager) studentByID(ctx context.Context, studentID int) (*coursepass.Student, error) {
	studentData, err := cm.repo.OneStudent(ctx, &db.StudentSearch{
		ID: &studentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed get student: %w", err)
	}
	if studentData == nil {
		return nil, coursepass.ErrStudentNotFound
	}

	return coursepass.NewStudent(studentData), nil
}

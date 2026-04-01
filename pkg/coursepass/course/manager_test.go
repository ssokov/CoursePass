package course

import (
	"testing"
	"time"

	"courses/pkg/coursepass"
	"courses/pkg/db"
	dbtest "courses/pkg/db/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type courseFixture struct {
	dbo     db.DB
	manager *Manager
	repo    db.CoursesRepo
}

func newCourseFixture(t *testing.T) courseFixture {
	t.Helper()

	dbo, logger := dbtest.Setup(t)

	return courseFixture{
		dbo:     dbo,
		manager: NewManager(dbo, logger),
		repo:    db.NewCoursesRepo(dbo),
	}
}

func TestCourseManager_Summary_AvailableOnly(t *testing.T) {
	// Arrange
	fx := newCourseFixture(t)

	now := time.Now()

	activeFrom := now.Add(-time.Hour)
	activeTo := now.Add(time.Hour)
	activeCourse, activeCleanup := dbtest.Course(
		t,
		fx.dbo.DB,
		&db.Course{
			AvailabilityType: "always",
			AvailableFrom:    &activeFrom,
			AvailableTo:      &activeTo,
			StatusID:         1,
		},
		dbtest.WithFakeCourse,
	)
	defer activeCleanup()

	futureFrom := now.Add(time.Hour)
	futureTo := now.Add(2 * time.Hour)
	futureCourse, futureCleanup := dbtest.Course(
		t,
		fx.dbo.DB,
		&db.Course{
			AvailabilityType: "always",
			AvailableFrom:    &futureFrom,
			AvailableTo:      &futureTo,
			StatusID:         1,
		},
		dbtest.WithFakeCourse,
	)
	defer futureCleanup()

	expiredFrom := now.Add(-2 * time.Hour)
	expiredTo := now.Add(-time.Hour)
	expiredCourse, expiredCleanup := dbtest.Course(
		t,
		fx.dbo.DB,
		&db.Course{
			AvailabilityType: "always",
			AvailableFrom:    &expiredFrom,
			AvailableTo:      &expiredTo,
			StatusID:         1,
		},
		dbtest.WithFakeCourse,
	)
	defer expiredCleanup()

	// Act
	list, err := fx.manager.List(t.Context(), 1, 50)

	// Assert
	require.NoError(t, err)
	assert.True(t, hasCourseID(list, activeCourse.ID))
	assert.False(t, hasCourseID(list, futureCourse.ID))
	assert.False(t, hasCourseID(list, expiredCourse.ID))
}

func TestCourseManager_ByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		fx := newCourseFixture(t)
		course, cleanup := dbtest.Course(
			t,
			fx.dbo.DB,
			&db.Course{
				AvailabilityType: "always",
				StatusID:         1,
			},
			dbtest.WithFakeCourse,
		)
		defer cleanup()

		// Act
		result, err := fx.manager.ByID(t.Context(), course.ID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, course.ID, result.ID)
		assert.Equal(t, course.Title, result.Title)
		assert.Equal(t, course.Description, result.Description)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		fx := newCourseFixture(t)

		// Act
		_, err := fx.manager.ByID(t.Context(), dbtest.NextID())

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, coursepass.ErrCourseNotFound)
	})
}

func TestCourseManager_Me(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		fx := newCourseFixture(t)
		login := "student_" + dbtest.NextStringID()
		email := "student_" + dbtest.NextStringID() + "@mail.test"
		student, cleanup := dbtest.Student(
			t,
			fx.dbo.DB,
			&db.Student{
				Login:    login,
				Email:    email,
				StatusID: 1,
			},
			dbtest.WithFakeStudent,
		)
		defer cleanup()

		// Act
		result, err := fx.manager.Me(t.Context(), student.ID)

		// Assert
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, student.ID, result.ID)
		assert.Equal(t, login, result.Login)
		assert.Equal(t, email, result.Email)
	})

	t.Run("not found", func(t *testing.T) {
		// Arrange
		fx := newCourseFixture(t)

		// Act
		_, err := fx.manager.Me(t.Context(), dbtest.NextID())

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, coursepass.ErrStudentNotFound)
	})
}

func hasCourseID(courses []coursepass.Course, courseID int) bool {
	for _, course := range courses {
		if course.ID == courseID {
			return true
		}
	}

	return false
}

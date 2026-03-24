package coursepass

import (
	"errors"
	"testing"
	"time"

	"courses/pkg/db"
	dbtest "courses/pkg/db/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type examFixture struct {
	dbo       db.DB
	manager   *ExamManager
	repo      db.CoursesRepo
	student   *db.Student
	course    *db.Course
	questions []*db.Question
	cleanup   func()
}

func newExamFixture(t *testing.T, questionCount int) examFixture {
	t.Helper()

	dbo, logger := dbtest.Setup(t)
	manager := NewExamManager(dbo, logger, "/media/")
	repo := db.NewCoursesRepo(dbo)

	var cleanups []func()

	student, studentCleanup := dbtest.Student(
		t,
		dbo.DB,
		&db.Student{StatusID: 1},
		dbtest.WithFakeStudent,
	)
	cleanups = append(cleanups, studentCleanup)

	now := time.Now()
	availableFrom := now.Add(-1 * time.Hour)
	availableTo := now.Add(1 * time.Hour)
	course, courseCleanup := dbtest.Course(
		t,
		dbo.DB,
		&db.Course{
			AvailabilityType: "always",
			AvailableFrom:    &availableFrom,
			AvailableTo:      &availableTo,
			StatusID:         1,
		},
		dbtest.WithFakeCourse,
	)
	cleanups = append(cleanups, courseCleanup)

	questions := make([]*db.Question, 0, questionCount)
	for i := 0; i < questionCount; i++ {
		question, questionCleanup := dbtest.Question(
			t,
			dbo.DB,
			&db.Question{
				CourseID:     course.ID,
				QuestionType: QuestionTypeSingleChoice,
				Options: db.QuestionOptions{
					{OptionID: 1, OptionText: "A", IsCorrect: true, DisplaySort: 1},
					{OptionID: 2, OptionText: "B", IsCorrect: false, DisplaySort: 2},
				},
			},
			dbtest.WithQuestionRelations,
			dbtest.WithFakeQuestion,
		)
		cleanups = append(cleanups, questionCleanup)
		questions = append(questions, question)
	}

	return examFixture{
		dbo:       dbo,
		manager:   manager,
		repo:      repo,
		student:   student,
		course:    course,
		questions: questions,
		cleanup: func() {
			for i := len(cleanups) - 1; i >= 0; i-- {
				cleanups[i]()
			}
		},
	}
}

func TestExamManager_Start(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		fx := newExamFixture(t, 2)
		defer fx.cleanup()

		// Act
		start, err := fx.manager.Start(t.Context(), fx.student.ID, fx.course.ID)

		// Assert
		require.NoError(t, err)
		require.Positive(t, start.ExamID)
		assert.Len(t, start.QuestionIDs, 2)
		assert.Nil(t, start.FinishedAt)

		defer func() {
			_, deleteErr := fx.repo.DeleteExam(t.Context(), start.ExamID)
			require.NoError(t, deleteErr)
		}()

		exam, findErr := fx.repo.OneExam(t.Context(), &db.ExamSearch{ID: &start.ExamID, StudentID: &fx.student.ID})
		require.NoError(t, findErr)
		require.NotNil(t, exam)
		assert.Equal(t, ExamStatusInProgress, exam.Status)
		assert.Equal(t, 2, len(exam.QuestionIDs))
	})

	t.Run("no questions", func(t *testing.T) {
		// Arrange
		fx := newExamFixture(t, 0)
		defer fx.cleanup()

		// Act
		_, err := fx.manager.Start(t.Context(), fx.student.ID, fx.course.ID)

		// Assert
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoQuestions))
	})
}

func TestExamManager_Question_NotInExam(t *testing.T) {
	// Arrange
	fx := newExamFixture(t, 1)
	defer fx.cleanup()

	start, err := fx.manager.Start(t.Context(), fx.student.ID, fx.course.ID)
	require.NoError(t, err)
	defer func() {
		_, deleteErr := fx.repo.DeleteExam(t.Context(), start.ExamID)
		require.NoError(t, deleteErr)
	}()

	// Act
	_, err = fx.manager.Question(t.Context(), fx.student.ID, dbtest.NextID(), start.ExamID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrQuestionNotInExam))
}

func TestExamManager_SaveAnswer_InvalidOptionIDs(t *testing.T) {
	// Arrange
	fx := newExamFixture(t, 1)
	defer fx.cleanup()

	start, err := fx.manager.Start(t.Context(), fx.student.ID, fx.course.ID)
	require.NoError(t, err)
	defer func() {
		_, deleteErr := fx.repo.DeleteExam(t.Context(), start.ExamID)
		require.NoError(t, deleteErr)
	}()

	// Act
	err = fx.manager.SaveAnswer(t.Context(), fx.student.ID, start.ExamID, fx.questions[0].ID, []int{999})

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidOptionIDs))
}

func TestExamManager_Submit_Success(t *testing.T) {
	// Arrange
	fx := newExamFixture(t, 2)
	defer fx.cleanup()

	start, err := fx.manager.Start(t.Context(), fx.student.ID, fx.course.ID)
	require.NoError(t, err)
	defer func() {
		_, deleteErr := fx.repo.DeleteExam(t.Context(), start.ExamID)
		require.NoError(t, deleteErr)
	}()

	require.NoError(t, fx.manager.SaveAnswer(t.Context(), fx.student.ID, start.ExamID, fx.questions[0].ID, []int{1}))
	require.NoError(t, fx.manager.SaveAnswer(t.Context(), fx.student.ID, start.ExamID, fx.questions[1].ID, []int{2}))

	// Act
	result, err := fx.manager.Submit(t.Context(), fx.student.ID, start.ExamID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, start.ExamID, result.ExamID)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 1, result.CorrectAnswers)
	assert.Equal(t, 50, result.FinalScore)
	assert.Equal(t, ExamStatusFailed, result.Status)

	exam, findErr := fx.repo.OneExam(t.Context(), &db.ExamSearch{ID: &start.ExamID, StudentID: &fx.student.ID})
	require.NoError(t, findErr)
	require.NotNil(t, exam)
	assert.Equal(t, ExamStatusFailed, exam.Status)
	require.NotNil(t, exam.FinishedAt)
}

func TestExamManager_MyList_OnlyFinished(t *testing.T) {
	// Arrange
	fx := newExamFixture(t, 1)
	defer fx.cleanup()

	inProgressExam, inProgressCleanup := dbtest.Exam(
		t,
		fx.dbo.DB,
		&db.Exam{
			CourseID:    fx.course.ID,
			StudentID:   fx.student.ID,
			QuestionIDs: []int{fx.questions[0].ID},
			Status:      ExamStatusInProgress,
			Answers:     db.ExamAnswers{},
		},
		dbtest.WithExamRelations,
		dbtest.WithFakeExam,
	)
	defer inProgressCleanup()

	finishedAt := time.Now()
	finalScore := 100.0
	finishedExam, finishedCleanup := dbtest.Exam(
		t,
		fx.dbo.DB,
		&db.Exam{
			CourseID:       fx.course.ID,
			StudentID:      fx.student.ID,
			QuestionIDs:    []int{fx.questions[0].ID},
			Status:         ExamStatusPassed,
			Answers:        db.ExamAnswers{},
			TotalQuestions: ptr(1),
			CorrectAnswers: ptr(1),
			FinalScore:     &finalScore,
			FinishedAt:     &finishedAt,
		},
		dbtest.WithExamRelations,
		dbtest.WithFakeExam,
	)
	defer finishedCleanup()

	// Act
	list, err := fx.manager.MyList(t.Context(), fx.student.ID, 1, 10)

	// Assert
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, finishedExam.ID, list[0].ExamID)
	assert.Contains(t, []string{ExamStatusPassed, ExamStatusFailed}, list[0].Status)
	assert.NotEqual(t, inProgressExam.ID, list[0].ExamID)
}

func ptr[T any](v T) *T {
	return &v
}

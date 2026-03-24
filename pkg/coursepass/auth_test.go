package coursepass

import (
	"errors"
	"testing"

	"courses/pkg/db"
	dbtest "courses/pkg/db/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type authFixture struct {
	dbo     db.DB
	manager *AuthManager
	repo    db.CoursesRepo
	authCfg AuthConfig
}

func newAuthFixture(t *testing.T) authFixture {
	t.Helper()

	dbo, logger := dbtest.Setup(t)
	authCfg := AuthConfig{
		JWTSecret:     "test-secret",
		JWTTTLSeconds: 3600,
	}

	return authFixture{
		dbo:     dbo,
		manager: NewAuthManager(dbo, logger, authCfg),
		repo:    db.NewCoursesRepo(dbo),
		authCfg: authCfg,
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	return string(hash)
}

func createStudent(
	t *testing.T,
	dbo db.DB,
	login, email, passwordHash string,
) (*db.Student, func()) {
	t.Helper()

	student, cleanup := dbtest.Student(
		t,
		dbo.DB,
		&db.Student{
			Login:        login,
			Email:        email,
			PasswordHash: passwordHash,
			StatusID:     1,
		},
		dbtest.WithFakeStudent,
	)

	return student, cleanup
}

func TestAuthManager_Register_Success(t *testing.T) {
	// Arrange
	fx := newAuthFixture(t)

	login := "student_" + dbtest.NextStringID()
	email := "student_" + dbtest.NextStringID() + "@mail.test"

	// Act
	token, err := fx.manager.Register(
		t.Context(),
		login,
		"password123",
		email,
		"John",
		"Doe",
	)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.Equal(t, fx.authCfg.JWTTTLSeconds, token.ExpiresIn)
	assert.Equal(t, "Bearer", token.TokenType)

	studentID, err := ValidateJWT(fx.authCfg, token.AccessToken)
	require.NoError(t, err)
	assert.Positive(t, studentID)

	student, err := fx.repo.OneStudent(t.Context(), &db.StudentSearch{ID: &studentID})
	require.NoError(t, err)
	require.NotNil(t, student)
	assert.Equal(t, login, student.Login)
	assert.Equal(t, email, student.Email)
	assert.Equal(t, "John", student.FirstName)
	assert.Equal(t, "Doe", student.LastName)
}

func TestAuthManager_Register_DuplicateLogin(t *testing.T) {
	// Arrange
	fx := newAuthFixture(t)

	passwordHash := hashPassword(t, "password123")

	existingLogin := "student_" + dbtest.NextStringID()
	_, cleanup := createStudent(
		t,
		fx.dbo,
		existingLogin,
		"student_"+dbtest.NextStringID()+"@mail.test",
		passwordHash,
	)
	defer cleanup()

	// Act
	_, err := fx.manager.Register(
		t.Context(),
		existingLogin,
		"password123",
		"student_"+dbtest.NextStringID()+"@mail.test",
		"John",
		"Doe",
	)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrLoginExists))
}

func TestAuthManager_Register_DuplicateEmail(t *testing.T) {
	// Arrange
	fx := newAuthFixture(t)

	passwordHash := hashPassword(t, "password123")

	existingEmail := "student_" + dbtest.NextStringID() + "@mail.test"
	_, cleanup := createStudent(
		t,
		fx.dbo,
		"student_"+dbtest.NextStringID(),
		existingEmail,
		passwordHash,
	)
	defer cleanup()

	// Act
	_, err := fx.manager.Register(
		t.Context(),
		"student_"+dbtest.NextStringID(),
		"password123",
		existingEmail,
		"John",
		"Doe",
	)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEmailExists))
}

func TestAuthManager_Login_Success(t *testing.T) {
	// Arrange
	fx := newAuthFixture(t)

	rawPassword := "password123"
	passwordHash := hashPassword(t, rawPassword)

	student, cleanup := createStudent(
		t,
		fx.dbo,
		"student_"+dbtest.NextStringID(),
		"student_"+dbtest.NextStringID()+"@mail.test",
		passwordHash,
	)
	defer cleanup()

	// Act
	token, err := fx.manager.Login(t.Context(), student.Login, rawPassword)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.Equal(t, "Bearer", token.TokenType)

	studentID, err := ValidateJWT(fx.authCfg, token.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, student.ID, studentID)
}

func TestAuthManager_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	fx := newAuthFixture(t)

	rawPassword := "password123"
	passwordHash := hashPassword(t, rawPassword)

	student, cleanup := createStudent(
		t,
		fx.dbo,
		"student_"+dbtest.NextStringID(),
		"student_"+dbtest.NextStringID()+"@mail.test",
		passwordHash,
	)
	defer cleanup()

	// Act
	_, err := fx.manager.Login(t.Context(), student.Login, "wrong-password")

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidCredentials))
}

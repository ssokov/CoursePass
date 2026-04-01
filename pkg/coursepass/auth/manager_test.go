package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"courses/pkg/coursepass"
	"courses/pkg/db"
	dbtest "courses/pkg/db/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type authFixture struct {
	dbo           db.DB
	manager       *Manager
	repo          db.CoursesRepo
	jwtSecret     string
	jwtTTLSeconds int
}

func newAuthFixture(t *testing.T) authFixture {
	t.Helper()

	dbo, logger := dbtest.Setup(t)
	jwtSecret := "test-secret"
	jwtTTLSeconds := 3600

	return authFixture{
		dbo:           dbo,
		manager:       NewManager(dbo, logger, jwtSecret, jwtTTLSeconds),
		repo:          db.NewCoursesRepo(dbo),
		jwtSecret:     jwtSecret,
		jwtTTLSeconds: jwtTTLSeconds,
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	return string(hash)
}

func TestHashPassword_DifferentPasswords(t *testing.T) {
	hashOne := hashPassword(t, "password123")
	hashTwo := hashPassword(t, "different-password-123")

	assert.NotEqual(t, hashOne, hashTwo)
}

func createStudent(t *testing.T, dbo db.DB, login, email, passwordHash string) (*db.Student, func()) {
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

func validateJWT(t *testing.T, secret, token string) int {
	t.Helper()

	parts := strings.Split(token, ".")
	require.Len(t, parts, 3)

	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(unsigned))
	require.NoError(t, err)

	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	require.Equal(t, expected, parts[2])

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	require.NoError(t, err)

	var claims struct {
		Sub string `json:"sub"`
	}
	require.NoError(t, json.Unmarshal(payloadRaw, &claims))

	studentID, err := strconv.Atoi(claims.Sub)
	require.NoError(t, err)

	return studentID
}

func TestAuthManager_Register_Success(t *testing.T) {
	// Arrange
	fx := newAuthFixture(t)

	login := "student_" + dbtest.NextStringID()
	email := "student_" + dbtest.NextStringID() + "@mail.test"

	// Act
	token, err := fx.manager.RegisterStudent(
		t.Context(),
		coursepass.StudentDraft{
			Login:     login,
			Password:  "password123",
			Email:     email,
			FirstName: "John",
			LastName:  "Doe",
		},
	)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.Equal(t, fx.jwtTTLSeconds, token.ExpiresIn)
	assert.Equal(t, "Bearer", token.TokenType)

	studentID := validateJWT(t, fx.jwtSecret, token.AccessToken)
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
	_, err := fx.manager.RegisterStudent(
		t.Context(),
		coursepass.StudentDraft{
			Login:     existingLogin,
			Password:  "password123",
			Email:     "student_" + dbtest.NextStringID() + "@mail.test",
			FirstName: "John",
			LastName:  "Doe",
		},
	)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, coursepass.ErrLoginExists)
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
	_, err := fx.manager.RegisterStudent(
		t.Context(),
		coursepass.StudentDraft{
			Login:     "student_" + dbtest.NextStringID(),
			Password:  "password123",
			Email:     existingEmail,
			FirstName: "John",
			LastName:  "Doe",
		},
	)

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, coursepass.ErrEmailExists)
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
	token, err := fx.manager.Login(t.Context(), coursepass.StudentLogin{
		Login:    student.Login,
		Password: rawPassword,
	})

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.Equal(t, "Bearer", token.TokenType)

	studentID := validateJWT(t, fx.jwtSecret, token.AccessToken)
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
	_, err := fx.manager.Login(t.Context(), coursepass.StudentLogin{
		Login:    student.Login,
		Password: "wrong-password",
	})

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, coursepass.ErrInvalidCredentials)
}

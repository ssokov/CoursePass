//nolint:dupl,funlen
package test

import (
	"testing"
	"time"

	"courses/pkg/db"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-pg/pg/v10/orm"
)

type UserOpFunc func(t *testing.T, dbo orm.DB, in *db.User) Cleaner

func User(t *testing.T, dbo orm.DB, in *db.User, ops ...UserOpFunc) (*db.User, Cleaner) {
	repo := db.NewCoursesRepo(dbo)
	var cleaners []Cleaner

	// Fill the incoming entity
	if in == nil {
		in = &db.User{}
	}

	// Check if PKs are provided
	if in.ID != 0 {
		// Fetch the entity by PK
		user, err := repo.UserByID(t.Context(), in.ID, repo.FullUser())
		if err != nil {
			t.Fatal(err)
		}

		// We must find the entity by PK
		if user == nil {
			t.Fatalf("the entity User is not found by provided PKs ID=%v", in.ID)
		}

		// Return if found without real cleanup
		return user, emptyClean
	}

	for _, op := range ops {
		if cl := op(t, dbo, in); cl != nil {
			cleaners = append(cleaners, cl)
		}
	}

	// Create the main entity
	user, err := repo.AddUser(t.Context(), in)
	if err != nil {
		t.Fatal(err)
	}

	return user, func() {
		if _, err := dbo.ModelContext(t.Context(), &db.User{ID: user.ID}).WherePK().Delete(); err != nil {
			t.Fatal(err)
		}
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithFakeUser(t *testing.T, dbo orm.DB, in *db.User) Cleaner {
	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	if in.Login == "" {
		in.Login = cutS(gofakeit.Word(), 64)
	}

	if in.Password == "" {
		in.Password = cutS(gofakeit.Password(true, true, true, false, false, 12), 64)
	}

	if in.AuthKey == "" {
		in.AuthKey = cutS(gofakeit.Sentence(3), 32)
	}

	if in.StatusID == 0 {
		in.StatusID = 1
	}

	return emptyClean
}

type CourseOpFunc func(t *testing.T, dbo orm.DB, in *db.Course) Cleaner

func Course(t *testing.T, dbo orm.DB, in *db.Course, ops ...CourseOpFunc) (*db.Course, Cleaner) {
	repo := db.NewCoursesRepo(dbo)
	var cleaners []Cleaner

	// Fill the incoming entity
	if in == nil {
		in = &db.Course{}
	}

	// Check if PKs are provided
	if in.ID != 0 {
		// Fetch the entity by PK
		course, err := repo.CourseByID(t.Context(), in.ID, repo.FullCourse())
		if err != nil {
			t.Fatal(err)
		}

		// We must find the entity by PK
		if course == nil {
			t.Fatalf("the entity Course is not found by provided PKs ID=%v", in.ID)
		}

		// Return if found without real cleanup
		return course, emptyClean
	}

	for _, op := range ops {
		if cl := op(t, dbo, in); cl != nil {
			cleaners = append(cleaners, cl)
		}
	}

	// Create the main entity
	course, err := repo.AddCourse(t.Context(), in)
	if err != nil {
		t.Fatal(err)
	}

	return course, func() {
		if _, err := dbo.ModelContext(t.Context(), &db.Course{ID: course.ID}).WherePK().Delete(); err != nil {
			t.Fatal(err)
		}
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithFakeCourse(t *testing.T, dbo orm.DB, in *db.Course) Cleaner {
	if in.Title == "" {
		in.Title = cutS(gofakeit.Sentence(10), 255)
	}

	if in.Description == "" {
		in.Description = cutS(gofakeit.Sentence(10), 0)
	}

	if in.AvailabilityType == "" {
		in.AvailabilityType = cutS(gofakeit.Sentence(10), 255)
	}

	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	if in.StatusID == 0 {
		in.StatusID = 1
	}

	return emptyClean
}

type ExamOpFunc func(t *testing.T, dbo orm.DB, in *db.Exam) Cleaner

func Exam(t *testing.T, dbo orm.DB, in *db.Exam, ops ...ExamOpFunc) (*db.Exam, Cleaner) {
	repo := db.NewCoursesRepo(dbo)
	var cleaners []Cleaner

	// Fill the incoming entity
	if in == nil {
		in = &db.Exam{}
	}

	// Check if PKs are provided
	if in.ID != 0 {
		// Fetch the entity by PK
		exam, err := repo.ExamByID(t.Context(), in.ID, repo.FullExam())
		if err != nil {
			t.Fatal(err)
		}

		// We must find the entity by PK
		if exam == nil {
			t.Fatalf("the entity Exam is not found by provided PKs ID=%v", in.ID)
		}

		// Return if found without real cleanup
		return exam, emptyClean
	}

	for _, op := range ops {
		if cl := op(t, dbo, in); cl != nil {
			cleaners = append(cleaners, cl)
		}
	}

	// Create the main entity
	exam, err := repo.AddExam(t.Context(), in)
	if err != nil {
		t.Fatal(err)
	}

	return exam, func() {
		if _, err := dbo.ModelContext(t.Context(), &db.Exam{ID: exam.ID}).WherePK().Delete(); err != nil {
			t.Fatal(err)
		}
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithExamRelations(t *testing.T, dbo orm.DB, in *db.Exam) Cleaner {
	var cleaners []Cleaner

	// Prepare main relations
	if in.Course == nil {
		in.Course = &db.Course{}
	}

	if in.QuestionIDs == nil {
		in.QuestionIDs = []int{}
	}

	if in.Student == nil {
		in.Student = &db.Student{}
	}

	// Check if all FKs are provided. Fill them into the main struct rels

	if in.CourseID != 0 {
		in.Course.ID = in.CourseID
	}

	if in.StudentID != 0 {
		in.Student.ID = in.StudentID
	}

	// Fetch the relation. It creates if the FKs are provided it fetch from DB by PKs. Else it creates new one.
	{
		for i := range in.QuestionIDs {
			_, relatedCleaner := Question(t, dbo, &db.Question{ID: in.QuestionIDs[i]}, WithQuestionRelations, WithFakeQuestion)

			cleaners = append(cleaners, relatedCleaner)
		}
	}

	// Fetch the relation. It creates if the FKs are provided it fetch from DB by PKs. Else it creates new one.
	{
		rel, relatedCleaner := Course(t, dbo, in.Course, WithFakeCourse)
		in.Course = rel
		in.CourseID = rel.ID

		cleaners = append(cleaners, relatedCleaner)
	}

	// Fetch the relation. It creates if the FKs are provided it fetch from DB by PKs. Else it creates new one.
	{
		rel, relatedCleaner := Student(t, dbo, in.Student, WithFakeStudent)
		in.Student = rel
		in.StudentID = rel.ID

		cleaners = append(cleaners, relatedCleaner)
	}

	return func() {
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithFakeExam(t *testing.T, dbo orm.DB, in *db.Exam) Cleaner {
	if in.Status == "" {
		in.Status = cutS(gofakeit.Sentence(10), 255)
	}

	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	return emptyClean
}

type QuestionOpFunc func(t *testing.T, dbo orm.DB, in *db.Question) Cleaner

func Question(t *testing.T, dbo orm.DB, in *db.Question, ops ...QuestionOpFunc) (*db.Question, Cleaner) {
	repo := db.NewCoursesRepo(dbo)
	var cleaners []Cleaner

	// Fill the incoming entity
	if in == nil {
		in = &db.Question{}
	}

	// Check if PKs are provided
	if in.ID != 0 {
		// Fetch the entity by PK
		question, err := repo.QuestionByID(t.Context(), in.ID, repo.FullQuestion())
		if err != nil {
			t.Fatal(err)
		}

		// We must find the entity by PK
		if question == nil {
			t.Fatalf("the entity Question is not found by provided PKs ID=%v", in.ID)
		}

		// Return if found without real cleanup
		return question, emptyClean
	}

	for _, op := range ops {
		if cl := op(t, dbo, in); cl != nil {
			cleaners = append(cleaners, cl)
		}
	}

	// Create the main entity
	question, err := repo.AddQuestion(t.Context(), in)
	if err != nil {
		t.Fatal(err)
	}

	return question, func() {
		if _, err := dbo.ModelContext(t.Context(), &db.Question{ID: question.ID}).WherePK().Delete(); err != nil {
			t.Fatal(err)
		}
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithQuestionRelations(t *testing.T, dbo orm.DB, in *db.Question) Cleaner {
	var cleaners []Cleaner

	// Prepare main relations
	if in.Course == nil {
		in.Course = &db.Course{}
	}

	// Check if all FKs are provided. Fill them into the main struct rels

	if in.CourseID != 0 {
		in.Course.ID = in.CourseID
	}

	// Fetch the relation. It creates if the FKs are provided it fetch from DB by PKs. Else it creates new one.
	{
		rel, relatedCleaner := Course(t, dbo, in.Course, WithFakeCourse)
		in.Course = rel
		in.CourseID = rel.ID

		cleaners = append(cleaners, relatedCleaner)
	}

	return func() {
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithFakeQuestion(t *testing.T, dbo orm.DB, in *db.Question) Cleaner {
	if in.QuestionText == "" {
		in.QuestionText = cutS(gofakeit.Sentence(10), 0)
	}

	if in.QuestionType == "" {
		in.QuestionType = cutS(gofakeit.Sentence(10), 255)
	}

	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	return emptyClean
}

type StudentOpFunc func(t *testing.T, dbo orm.DB, in *db.Student) Cleaner

func Student(t *testing.T, dbo orm.DB, in *db.Student, ops ...StudentOpFunc) (*db.Student, Cleaner) {
	repo := db.NewCoursesRepo(dbo)
	var cleaners []Cleaner

	// Fill the incoming entity
	if in == nil {
		in = &db.Student{}
	}

	// Check if PKs are provided
	if in.ID != 0 {
		// Fetch the entity by PK
		student, err := repo.StudentByID(t.Context(), in.ID, repo.FullStudent())
		if err != nil {
			t.Fatal(err)
		}

		// We must find the entity by PK
		if student == nil {
			t.Fatalf("the entity Student is not found by provided PKs ID=%v", in.ID)
		}

		// Return if found without real cleanup
		return student, emptyClean
	}

	for _, op := range ops {
		if cl := op(t, dbo, in); cl != nil {
			cleaners = append(cleaners, cl)
		}
	}

	// Create the main entity
	student, err := repo.AddStudent(t.Context(), in)
	if err != nil {
		t.Fatal(err)
	}

	return student, func() {
		if _, err := dbo.ModelContext(t.Context(), &db.Student{ID: student.ID}).WherePK().Delete(); err != nil {
			t.Fatal(err)
		}
		// Clean up related entities from the last to the first
		for i := len(cleaners) - 1; i >= 0; i-- {
			cleaners[i]()
		}
	}
}

func WithFakeStudent(t *testing.T, dbo orm.DB, in *db.Student) Cleaner {
	if in.Login == "" {
		in.Login = cutS(gofakeit.Word(), 255)
	}

	if in.PasswordHash == "" {
		in.PasswordHash = cutS(gofakeit.Sentence(10), 255)
	}

	if in.FirstName == "" {
		in.FirstName = cutS(gofakeit.Sentence(10), 255)
	}

	if in.LastName == "" {
		in.LastName = cutS(gofakeit.Sentence(10), 255)
	}

	if in.Email == "" {
		in.Email = cutS(gofakeit.Email(), 255)
	}

	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	if in.StatusID == 0 {
		in.StatusID = 1
	}

	return emptyClean
}

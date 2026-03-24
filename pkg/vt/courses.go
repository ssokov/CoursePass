package vt

import (
	"context"

	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type CourseService struct {
	zenrpc.Service
	embedlog.Logger
	coursesRepo db.CoursesRepo
}

func NewCourseService(dbo db.DB, logger embedlog.Logger) *CourseService {
	return &CourseService{
		Logger:      logger,
		coursesRepo: db.NewCoursesRepo(dbo),
	}
}

func (s CourseService) dbSort(ops *ViewOps) db.OpFunc {
	v := s.coursesRepo.DefaultCourseSort()
	if ops == nil {
		return v
	}

	switch ops.SortColumn {
	case db.Columns.Course.ID, db.Columns.Course.Title, db.Columns.Course.Description, db.Columns.Course.AvailabilityType, db.Columns.Course.AvailableFrom, db.Columns.Course.AvailableTo, db.Columns.Course.TimeLimitMinutes, db.Columns.Course.CreatedAt, db.Columns.Course.StatusID:
		v = db.WithSort(db.NewSortField(ops.SortColumn, ops.SortDesc))
	}

	return v
}

// Count returns count Courses according to conditions in search params.
//
//zenrpc:search CourseSearch
//zenrpc:return int
//zenrpc:500 Internal Error
func (s CourseService) Count(ctx context.Context, search *CourseSearch) (int, error) {
	count, err := s.coursesRepo.CountCourses(ctx, search.ToDB())
	if err != nil {
		return 0, InternalError(err)
	}
	return count, nil
}

// Get returns а list of Courses according to conditions in search params.
//
//zenrpc:search CourseSearch
//zenrpc:viewOps ViewOps
//zenrpc:return []CourseSummary
//zenrpc:500 Internal Error
func (s CourseService) Get(ctx context.Context, search *CourseSearch, viewOps *ViewOps) ([]CourseSummary, error) {
	list, err := s.coursesRepo.CoursesByFilters(ctx, search.ToDB(), viewOps.Pager(), s.dbSort(viewOps), s.coursesRepo.FullCourse())
	if err != nil {
		return nil, InternalError(err)
	}
	courses := make([]CourseSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		if course := NewCourseSummary(&list[i]); course != nil {
			courses = append(courses, *course)
		}
	}
	return courses, nil
}

// GetByID returns a Course by its ID.
//
//zenrpc:id int
//zenrpc:return Course
//zenrpc:500 Internal Error
//zenrpc:404 Not Found
func (s CourseService) GetByID(ctx context.Context, id int) (*Course, error) {
	db, err := s.byID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewCourse(db), nil
}

func (s CourseService) byID(ctx context.Context, id int) (*db.Course, error) {
	db, err := s.coursesRepo.CourseByID(ctx, id, s.coursesRepo.FullCourse())
	if err != nil {
		return nil, InternalError(err)
	} else if db == nil {
		return nil, ErrNotFound
	}
	return db, nil
}

// Add adds a Course from the query.
//
//zenrpc:course Course
//zenrpc:return Course
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
func (s CourseService) Add(ctx context.Context, course Course) (*Course, error) {
	if ve := s.isValid(ctx, course, false); ve.HasErrors() {
		return nil, ve.Error()
	}

	db, err := s.coursesRepo.AddCourse(ctx, course.ToDB())
	if err != nil {
		return nil, InternalError(err)
	}
	return NewCourse(db), nil
}

// Update updates the Course data identified by id from the query.
//
//zenrpc:courses Course
//zenrpc:return Course
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s CourseService) Update(ctx context.Context, course Course) (bool, error) {
	if _, err := s.byID(ctx, course.ID); err != nil {
		return false, err
	}

	if ve := s.isValid(ctx, course, true); ve.HasErrors() {
		return false, ve.Error()
	}

	ok, err := s.coursesRepo.UpdateCourse(ctx, course.ToDB())
	if err != nil {
		return false, InternalError(err)
	}
	return ok, nil
}

// Delete deletes the Course by its ID.
//
//zenrpc:id int
//zenrpc:return isDeleted
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s CourseService) Delete(ctx context.Context, id int) (bool, error) {
	if _, err := s.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := s.coursesRepo.DeleteCourse(ctx, id)
	if err != nil {
		return false, InternalError(err)
	}
	return ok, err
}

// Validate verifies that Course data is valid.
//
//zenrpc:course Course
//zenrpc:return []FieldError
//zenrpc:500 Internal Error
func (s CourseService) Validate(ctx context.Context, course Course) ([]FieldError, error) {
	isUpdate := course.ID != 0
	if isUpdate {
		_, err := s.byID(ctx, course.ID)
		if err != nil {
			return nil, err
		}
	}

	ve := s.isValid(ctx, course, isUpdate)
	if ve.HasInternalError() {
		return nil, ve.Error()
	}

	return ve.Fields(), nil
}

func (s CourseService) isValid(ctx context.Context, course Course, isUpdate bool) Validator {
	var v Validator

	if v.CheckBasic(ctx, course); v.HasInternalError() {
		return v
	}

	// custom validation starts here
	return v
}

type ExamService struct {
	zenrpc.Service
	embedlog.Logger
	coursesRepo db.CoursesRepo
}

func NewExamService(dbo db.DB, logger embedlog.Logger) *ExamService {
	return &ExamService{
		Logger:      logger,
		coursesRepo: db.NewCoursesRepo(dbo),
	}
}

func (s ExamService) dbSort(ops *ViewOps) db.OpFunc {
	v := s.coursesRepo.DefaultExamSort()
	if ops == nil {
		return v
	}

	switch ops.SortColumn {
	case db.Columns.Exam.ID, db.Columns.Exam.CourseID, db.Columns.Exam.StudentID, db.Columns.Exam.TotalQuestions, db.Columns.Exam.CorrectAnswers, db.Columns.Exam.Status, db.Columns.Exam.FinalScore, db.Columns.Exam.FinishedAt, db.Columns.Exam.CreatedAt:
		v = db.WithSort(db.NewSortField(ops.SortColumn, ops.SortDesc))
	}

	return v
}

// Count returns count Exams according to conditions in search params.
//
//zenrpc:search ExamSearch
//zenrpc:return int
//zenrpc:500 Internal Error
func (s ExamService) Count(ctx context.Context, search *ExamSearch) (int, error) {
	count, err := s.coursesRepo.CountExams(ctx, search.ToDB())
	if err != nil {
		return 0, InternalError(err)
	}
	return count, nil
}

// Get returns а list of Exams according to conditions in search params.
//
//zenrpc:search ExamSearch
//zenrpc:viewOps ViewOps
//zenrpc:return []ExamSummary
//zenrpc:500 Internal Error
func (s ExamService) Get(ctx context.Context, search *ExamSearch, viewOps *ViewOps) ([]ExamSummary, error) {
	list, err := s.coursesRepo.ExamsByFilters(ctx, search.ToDB(), viewOps.Pager(), s.dbSort(viewOps), s.coursesRepo.FullExam())
	if err != nil {
		return nil, InternalError(err)
	}
	exams := make([]ExamSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		if exam := NewExamSummary(&list[i]); exam != nil {
			exams = append(exams, *exam)
		}
	}
	return exams, nil
}

// GetByID returns a Exam by its ID.
//
//zenrpc:id int
//zenrpc:return Exam
//zenrpc:500 Internal Error
//zenrpc:404 Not Found
func (s ExamService) GetByID(ctx context.Context, id int) (*Exam, error) {
	db, err := s.byID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewExam(db), nil
}

func (s ExamService) byID(ctx context.Context, id int) (*db.Exam, error) {
	db, err := s.coursesRepo.ExamByID(ctx, id, s.coursesRepo.FullExam())
	if err != nil {
		return nil, InternalError(err)
	} else if db == nil {
		return nil, ErrNotFound
	}
	return db, nil
}

// Add adds a Exam from the query.
//
//zenrpc:exam Exam
//zenrpc:return Exam
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
func (s ExamService) Add(ctx context.Context, exam Exam) (*Exam, error) {
	if ve := s.isValid(ctx, exam, false); ve.HasErrors() {
		return nil, ve.Error()
	}

	db, err := s.coursesRepo.AddExam(ctx, exam.ToDB())
	if err != nil {
		return nil, InternalError(err)
	}
	return NewExam(db), nil
}

// Update updates the Exam data identified by id from the query.
//
//zenrpc:exams Exam
//zenrpc:return Exam
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s ExamService) Update(ctx context.Context, exam Exam) (bool, error) {
	if _, err := s.byID(ctx, exam.ID); err != nil {
		return false, err
	}

	if ve := s.isValid(ctx, exam, true); ve.HasErrors() {
		return false, ve.Error()
	}

	ok, err := s.coursesRepo.UpdateExam(ctx, exam.ToDB())
	if err != nil {
		return false, InternalError(err)
	}
	return ok, nil
}

// Delete deletes the Exam by its ID.
//
//zenrpc:id int
//zenrpc:return isDeleted
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s ExamService) Delete(ctx context.Context, id int) (bool, error) {
	if _, err := s.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := s.coursesRepo.DeleteExam(ctx, id)
	if err != nil {
		return false, InternalError(err)
	}
	return ok, err
}

// Validate verifies that Exam data is valid.
//
//zenrpc:exam Exam
//zenrpc:return []FieldError
//zenrpc:500 Internal Error
func (s ExamService) Validate(ctx context.Context, exam Exam) ([]FieldError, error) {
	isUpdate := exam.ID != 0
	if isUpdate {
		_, err := s.byID(ctx, exam.ID)
		if err != nil {
			return nil, err
		}
	}

	ve := s.isValid(ctx, exam, isUpdate)
	if ve.HasInternalError() {
		return nil, ve.Error()
	}

	return ve.Fields(), nil
}

func (s ExamService) isValid(ctx context.Context, exam Exam, isUpdate bool) Validator {
	var v Validator

	if v.CheckBasic(ctx, exam); v.HasInternalError() {
		return v
	}

	// check fks
	if exam.CourseID != 0 {
		item, err := s.coursesRepo.CourseByID(ctx, exam.CourseID)
		if err != nil {
			v.SetInternalError(err)
		} else if item == nil {
			v.Append("courseId", FieldErrorIncorrect)
		}
	}

	if exam.StudentID != 0 {
		item, err := s.coursesRepo.StudentByID(ctx, exam.StudentID)
		if err != nil {
			v.SetInternalError(err)
		} else if item == nil {
			v.Append("studentId", FieldErrorIncorrect)
		}
	}

	if len(exam.QuestionIDs) != 0 {
		items, err := s.coursesRepo.QuestionsByFilters(ctx, &db.QuestionSearch{IDs: exam.QuestionIDs}, db.PagerNoLimit)
		if err != nil {
			v.SetInternalError(err)
		} else if len(items) != len(exam.QuestionIDs) {
			v.Append("questionIds", FieldErrorIncorrect)
		}
	}
	// custom validation starts here
	return v
}

type QuestionService struct {
	zenrpc.Service
	embedlog.Logger
	coursesRepo db.CoursesRepo
	vfsRepo     db.VfsRepo
}

func NewQuestionService(dbo db.DB, logger embedlog.Logger) *QuestionService {
	return &QuestionService{
		Logger:      logger,
		coursesRepo: db.NewCoursesRepo(dbo),
		vfsRepo:     db.NewVfsRepo(dbo),
	}
}

func (s QuestionService) dbSort(ops *ViewOps) db.OpFunc {
	v := s.coursesRepo.DefaultQuestionSort()
	if ops == nil {
		return v
	}

	switch ops.SortColumn {
	case db.Columns.Question.ID, db.Columns.Question.CourseID, db.Columns.Question.PhotoFileID, db.Columns.Question.QuestionText, db.Columns.Question.QuestionType, db.Columns.Question.CreatedAt:
		v = db.WithSort(db.NewSortField(ops.SortColumn, ops.SortDesc))
	}

	return v
}

// Count returns count Questions according to conditions in search params.
//
//zenrpc:search QuestionSearch
//zenrpc:return int
//zenrpc:500 Internal Error
func (s QuestionService) Count(ctx context.Context, search *QuestionSearch) (int, error) {
	count, err := s.coursesRepo.CountQuestions(ctx, search.ToDB())
	if err != nil {
		return 0, InternalError(err)
	}
	return count, nil
}

// Get returns а list of Questions according to conditions in search params.
//
//zenrpc:search QuestionSearch
//zenrpc:viewOps ViewOps
//zenrpc:return []QuestionSummary
//zenrpc:500 Internal Error
func (s QuestionService) Get(ctx context.Context, search *QuestionSearch, viewOps *ViewOps) ([]QuestionSummary, error) {
	list, err := s.coursesRepo.QuestionsByFilters(ctx, search.ToDB(), viewOps.Pager(), s.dbSort(viewOps), s.coursesRepo.FullQuestion())
	if err != nil {
		return nil, InternalError(err)
	}
	questions := make([]QuestionSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		if question := NewQuestionSummary(&list[i]); question != nil {
			questions = append(questions, *question)
		}
	}
	return questions, nil
}

// GetByID returns a Question by its ID.
//
//zenrpc:id int
//zenrpc:return Question
//zenrpc:500 Internal Error
//zenrpc:404 Not Found
func (s QuestionService) GetByID(ctx context.Context, id int) (*Question, error) {
	db, err := s.byID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewQuestion(db), nil
}

func (s QuestionService) byID(ctx context.Context, id int) (*db.Question, error) {
	db, err := s.coursesRepo.QuestionByID(ctx, id, s.coursesRepo.FullQuestion())
	if err != nil {
		return nil, InternalError(err)
	} else if db == nil {
		return nil, ErrNotFound
	}
	return db, nil
}

// Add adds a Question from the query.
//
//zenrpc:question Question
//zenrpc:return Question
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
func (s QuestionService) Add(ctx context.Context, question Question) (*Question, error) {
	if ve := s.isValid(ctx, question, false); ve.HasErrors() {
		return nil, ve.Error()
	}

	db, err := s.coursesRepo.AddQuestion(ctx, question.ToDB())
	if err != nil {
		return nil, InternalError(err)
	}
	return NewQuestion(db), nil
}

// Update updates the Question data identified by id from the query.
//
//zenrpc:questions Question
//zenrpc:return Question
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s QuestionService) Update(ctx context.Context, question Question) (bool, error) {
	if _, err := s.byID(ctx, question.ID); err != nil {
		return false, err
	}

	if ve := s.isValid(ctx, question, true); ve.HasErrors() {
		return false, ve.Error()
	}

	ok, err := s.coursesRepo.UpdateQuestion(ctx, question.ToDB())
	if err != nil {
		return false, InternalError(err)
	}
	return ok, nil
}

// Delete deletes the Question by its ID.
//
//zenrpc:id int
//zenrpc:return isDeleted
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s QuestionService) Delete(ctx context.Context, id int) (bool, error) {
	if _, err := s.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := s.coursesRepo.DeleteQuestion(ctx, id)
	if err != nil {
		return false, InternalError(err)
	}
	return ok, err
}

// Validate verifies that Question data is valid.
//
//zenrpc:question Question
//zenrpc:return []FieldError
//zenrpc:500 Internal Error
func (s QuestionService) Validate(ctx context.Context, question Question) ([]FieldError, error) {
	isUpdate := question.ID != 0
	if isUpdate {
		_, err := s.byID(ctx, question.ID)
		if err != nil {
			return nil, err
		}
	}

	ve := s.isValid(ctx, question, isUpdate)
	if ve.HasInternalError() {
		return nil, ve.Error()
	}

	return ve.Fields(), nil
}

func (s QuestionService) isValid(ctx context.Context, question Question, isUpdate bool) Validator {
	var v Validator

	if v.CheckBasic(ctx, question); v.HasInternalError() {
		return v
	}

	// check fks
	if question.CourseID != 0 {
		item, err := s.coursesRepo.CourseByID(ctx, question.CourseID)
		if err != nil {
			v.SetInternalError(err)
		} else if item == nil {
			v.Append("courseId", FieldErrorIncorrect)
		}
	}

	if question.PhotoFileID != nil {
		item, err := s.vfsRepo.VfsFileByID(ctx, *question.PhotoFileID)
		if err != nil {
			v.SetInternalError(err)
		} else if item == nil {
			v.Append("photoFileId", FieldErrorIncorrect)
		}
	}

	// custom validation starts here
	return v
}

type StudentService struct {
	zenrpc.Service
	embedlog.Logger
	coursesRepo db.CoursesRepo
}

func NewStudentService(dbo db.DB, logger embedlog.Logger) *StudentService {
	return &StudentService{
		Logger:      logger,
		coursesRepo: db.NewCoursesRepo(dbo),
	}
}

func (s StudentService) dbSort(ops *ViewOps) db.OpFunc {
	v := s.coursesRepo.DefaultStudentSort()
	if ops == nil {
		return v
	}

	switch ops.SortColumn {
	case db.Columns.Student.ID, db.Columns.Student.Login, db.Columns.Student.PasswordHash, db.Columns.Student.FirstName, db.Columns.Student.LastName, db.Columns.Student.Email, db.Columns.Student.CreatedAt, db.Columns.Student.StatusID:
		v = db.WithSort(db.NewSortField(ops.SortColumn, ops.SortDesc))
	}

	return v
}

// Count returns count Students according to conditions in search params.
//
//zenrpc:search StudentSearch
//zenrpc:return int
//zenrpc:500 Internal Error
func (s StudentService) Count(ctx context.Context, search *StudentSearch) (int, error) {
	count, err := s.coursesRepo.CountStudents(ctx, search.ToDB())
	if err != nil {
		return 0, InternalError(err)
	}
	return count, nil
}

// Get returns а list of Students according to conditions in search params.
//
//zenrpc:search StudentSearch
//zenrpc:viewOps ViewOps
//zenrpc:return []StudentSummary
//zenrpc:500 Internal Error
func (s StudentService) Get(ctx context.Context, search *StudentSearch, viewOps *ViewOps) ([]StudentSummary, error) {
	list, err := s.coursesRepo.StudentsByFilters(ctx, search.ToDB(), viewOps.Pager(), s.dbSort(viewOps), s.coursesRepo.FullStudent())
	if err != nil {
		return nil, InternalError(err)
	}
	students := make([]StudentSummary, 0, len(list))
	for i := 0; i < len(list); i++ {
		if student := NewStudentSummary(&list[i]); student != nil {
			students = append(students, *student)
		}
	}
	return students, nil
}

// GetByID returns a Student by its ID.
//
//zenrpc:id int
//zenrpc:return Student
//zenrpc:500 Internal Error
//zenrpc:404 Not Found
func (s StudentService) GetByID(ctx context.Context, id int) (*Student, error) {
	db, err := s.byID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewStudent(db), nil
}

func (s StudentService) byID(ctx context.Context, id int) (*db.Student, error) {
	db, err := s.coursesRepo.StudentByID(ctx, id, s.coursesRepo.FullStudent())
	if err != nil {
		return nil, InternalError(err)
	} else if db == nil {
		return nil, ErrNotFound
	}
	return db, nil
}

// Add adds a Student from the query.
//
//zenrpc:student Student
//zenrpc:return Student
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
func (s StudentService) Add(ctx context.Context, student Student) (*Student, error) {
	if ve := s.isValid(ctx, student, false); ve.HasErrors() {
		return nil, ve.Error()
	}

	db, err := s.coursesRepo.AddStudent(ctx, student.ToDB())
	if err != nil {
		return nil, InternalError(err)
	}
	return NewStudent(db), nil
}

// Update updates the Student data identified by id from the query.
//
//zenrpc:students Student
//zenrpc:return Student
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s StudentService) Update(ctx context.Context, student Student) (bool, error) {
	if _, err := s.byID(ctx, student.ID); err != nil {
		return false, err
	}

	if ve := s.isValid(ctx, student, true); ve.HasErrors() {
		return false, ve.Error()
	}

	ok, err := s.coursesRepo.UpdateStudent(ctx, student.ToDB())
	if err != nil {
		return false, InternalError(err)
	}
	return ok, nil
}

// Delete deletes the Student by its ID.
//
//zenrpc:id int
//zenrpc:return isDeleted
//zenrpc:500 Internal Error
//zenrpc:400 Validation Error
//zenrpc:404 Not Found
func (s StudentService) Delete(ctx context.Context, id int) (bool, error) {
	if _, err := s.byID(ctx, id); err != nil {
		return false, err
	}

	ok, err := s.coursesRepo.DeleteStudent(ctx, id)
	if err != nil {
		return false, InternalError(err)
	}
	return ok, err
}

// Validate verifies that Student data is valid.
//
//zenrpc:student Student
//zenrpc:return []FieldError
//zenrpc:500 Internal Error
func (s StudentService) Validate(ctx context.Context, student Student) ([]FieldError, error) {
	isUpdate := student.ID != 0
	if isUpdate {
		_, err := s.byID(ctx, student.ID)
		if err != nil {
			return nil, err
		}
	}

	ve := s.isValid(ctx, student, isUpdate)
	if ve.HasInternalError() {
		return nil, ve.Error()
	}

	return ve.Fields(), nil
}

func (s StudentService) isValid(ctx context.Context, student Student, isUpdate bool) Validator {
	var v Validator

	if v.CheckBasic(ctx, student); v.HasInternalError() {
		return v
	}

	// custom validation starts here
	return v
}

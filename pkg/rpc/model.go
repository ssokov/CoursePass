package rpc

type RegisterRequest struct {
	Login     string `json:"login" validate:"required,min=3,max=64"`
	Password  string `json:"password" validate:"required,min=6,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	FirstName string `json:"firstName" validate:"required,max=255"`
	LastName  string `json:"lastName" validate:"required,max=255"`
}

type RegisterResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,max=64"`
	Password string `json:"password" validate:"required,max=255"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type MeResponse struct {
	StudentID int    `json:"studentId"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Course struct {
	CourseID      int     `json:"courseId"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	TimeLimit     *int    `json:"timeLimit"`
	AvailableType string  `json:"availableType"`
	AvailableFrom *string `json:"availableFrom"`
	AvailableTo   *string `json:"availableTo"`
}

type ByIDResponse struct {
	Course Course `json:"coursepass"`
}

type ByIDRequest struct {
	CourseID int `json:"courseId"`
}

type CourseSummary struct {
	CourseID      int     `json:"courseId"`
	Title         string  `json:"title"`
	TimeLimit     *int    `json:"timeLimit"`
	AvailableType string  `json:"availableType"`
	AvailableFrom *string `json:"availableFrom"`
	AvailableTo   *string `json:"availableTo"`
}

type ListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type ListResponse struct {
	Courses []CourseSummary `json:"courses"`
}

type ExamStartRequest struct {
	CourseID int `json:"courseId"`
}

type ExamQuestionRequest struct {
	ExamID     int `json:"examId"`
	QuestionID int `json:"questionId"`
}

type ExamSaveAnswerRequest struct {
	ExamID     int   `json:"examId"`
	QuestionID int   `json:"questionId"`
	OptionIDs  []int `json:"optionIds"`
}

type ExamSubmitRequest struct {
	ExamID int `json:"examId"`
}

type ExamMyListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type ExamStartResponse struct {
	ExamID      int     `json:"examId"`
	QuestionIDs []int   `json:"questionIds"`
	StartedAt   string  `json:"startedAt"`
	FinishedAt  *string `json:"finishedAt"`
}

type Question struct {
	QuestionID   int              `json:"questionId"`
	QuestionText string           `json:"questionText"`
	QuestionType string           `json:"questionType"`
	PhotoURL     *string          `json:"photoUrl"`
	Options      []QuestionOption `json:"options"`
}

type QuestionOption struct {
	OptionID   int    `json:"optionId"`
	OptionText string `json:"optionText"`
}

type ExamResult struct {
	ExamID         int    `json:"examId"`
	Status         string `json:"status"`
	FinalScore     int    `json:"finalScore"`
	CorrectAnswers int    `json:"correctAnswers"`
	TotalQuestions int    `json:"totalQuestions"`
}

type ExamSummary struct {
	ExamID     int    `json:"examId"`
	CourseID   int    `json:"courseId"`
	Status     string `json:"status"`
	FinalScore int    `json:"finalScore"`
	FinishedAt string `json:"finishedAt"`
}

type ExamMyListResponse struct {
	Exams []ExamSummary `json:"exams"`
}

type SaveAnswerRequest struct {
	ExamID     int   `json:"examId"`
	QuestionID int   `json:"questionId"`
	OptionIDs  []int `json:"optionIds"`
}

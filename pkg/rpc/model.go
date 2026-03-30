package rpc

type Token struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
	TokenType   string `json:"tokenType"`
}

type Student struct {
	StudentID int    `json:"studentId"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RegisterStudent struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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

type CourseSummary struct {
	CourseID      int     `json:"courseId"`
	Title         string  `json:"title"`
	TimeLimit     *int    `json:"timeLimit"`
	AvailableType string  `json:"availableType"`
	AvailableFrom *string `json:"availableFrom"`
	AvailableTo   *string `json:"availableTo"`
}

type ExamStart struct {
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

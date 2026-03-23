package db

type ExamAnswer struct {
	QuestionID int   `json:"questionId"`
	OptionIDs  []int `json:"optionIds"`
}

type ExamAnswers []ExamAnswer

type QuestionOption struct {
	OptionID    int    `json:"optionId"`
	OptionText  string `json:"optionText"`
	IsCorrect   bool   `json:"isCorrect,omitempty"`
	DisplaySort int    `json:"displaySort,omitempty"`
}

type QuestionOptions []QuestionOption

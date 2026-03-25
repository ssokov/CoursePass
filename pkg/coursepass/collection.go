package coursepass

//go:generate colgen
//colgen:Question
//colgen:Question:QuestionID,Index(QuestionID)
//colgen:QuestionOption
//colgen:QuestionOption:OptionID,Index(OptionID),Group(IsCorrect)
//colgen:ExamAnswer
//colgen:ExamAnswer:Index(QuestionID)

func Map[S, T any](in []S, convert func(S) T) []T {
	out := make([]T, len(in))
	for i := range in {
		out[i] = convert(in[i])
	}

	return out
}

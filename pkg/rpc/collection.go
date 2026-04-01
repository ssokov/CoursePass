package rpc

//go:generate colgen -imports courses/pkg/coursepass
//colgen:CourseSummary:map(coursepass.Course)
//colgen:ExamSummary:map(coursepass.Exam)
//colgen:QuestionOption:Map(coursepass.QuestionOption)
//colgen:FieldError:map(coursepass.FieldError)

func Map[S, T any](in []S, convert func(S) T) []T {
	out := make([]T, len(in))
	for i := range in {
		out[i] = convert(in[i])
	}

	return out
}

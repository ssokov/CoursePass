package rpc

//go:generate colgen -imports courses/pkg/coursepass
//colgen:CourseSummary:map(coursepass.CourseSummary)
//colgen:QuestionOption:Map(coursepass.QuestionOption)

func Map[S, T any](in []S, convert func(S) T) []T {
	out := make([]T, len(in))
	for i := range in {
		out[i] = convert(in[i])
	}

	return out
}

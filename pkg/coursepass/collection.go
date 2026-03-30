package coursepass

import "courses/pkg/db"

//go:generate colgen -imports courses/pkg/db
//colgen:Question
//colgen:Course:MapP(db.Course)
//colgen:Exam:MapP(db.Exam)
//colgen:Question:MapP(db.Question)

func MapP[T, M any](in []T, convert func(*T) *M) []M {
	out := make([]M, len(in))
	for i := range in {
		out[i] = *convert(&in[i])
	}

	return out
}

type ExamAnswers db.ExamAnswers

func (ll ExamAnswers) IndexByQuestionID() map[int]db.ExamAnswer {
	r := make(map[int]db.ExamAnswer, len(ll))
	for i := range ll {
		r[ll[i].QuestionID] = ll[i]
	}
	return r
}

type QuestionOptions db.QuestionOptions

func (ll QuestionOptions) OptionIDs() []int {
	r := make([]int, len(ll))
	for i := range ll {
		r[i] = ll[i].OptionID
	}
	return r
}

func (ll QuestionOptions) IndexByOptionID() map[int]db.QuestionOption {
	r := make(map[int]db.QuestionOption, len(ll))
	for i := range ll {
		r[ll[i].OptionID] = ll[i]
	}
	return r
}

func (ll QuestionOptions) GroupByIsCorrect() map[bool]QuestionOptions {
	r := make(map[bool]QuestionOptions, len(ll))
	for i := range ll {
		r[ll[i].IsCorrect] = append(r[ll[i].IsCorrect], ll[i])
	}
	return r
}

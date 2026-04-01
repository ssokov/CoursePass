package coursepass

import "courses/pkg/db"

//go:generate colgen -imports courses/pkg/db
//colgen:Question
//colgen:Course:MapP(db.Course)
//colgen:Exam:MapP(db.Exam)
//colgen:Question:MapP(db.Question)
//colgen:QuestionOption:MapP(db.QuestionOption)

func MapP[T, M any](in []T, convert func(*T) *M) []M {
	out := make([]M, len(in))
	for i := range in {
		out[i] = *convert(&in[i])
	}

	return out
}

type ExamAnswers db.ExamAnswers

func (ea ExamAnswers) IndexByQuestionID() map[int]db.ExamAnswer {
	r := make(map[int]db.ExamAnswer, len(ea))
	for i := range ea {
		r[ea[i].QuestionID] = ea[i]
	}
	return r
}

type QuestionOptions []QuestionOption

func (qo QuestionOptions) OptionIDs() []int {
	r := make([]int, len(qo))
	for i := range qo {
		r[i] = qo[i].OptionID
	}
	return r
}

func (qo QuestionOptions) IndexByOptionID() map[int]QuestionOption {
	r := make(map[int]QuestionOption, len(qo))
	for i := range qo {
		r[qo[i].OptionID] = qo[i]
	}
	return r
}

func (qo QuestionOptions) GroupByIsCorrect() map[bool]QuestionOptions {
	r := make(map[bool]QuestionOptions, len(qo))
	for i := range qo {
		r[qo[i].IsCorrect] = append(r[qo[i].IsCorrect], qo[i])
	}
	return r
}

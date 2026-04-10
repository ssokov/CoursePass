/* Code generated from jsonrpc schema by rpcgen v2.5.x with typescript v1.0.0; DO NOT EDIT. */
/* eslint-disable */
export interface IAuthLoginParams {
  studentLogin: IStudentLogin
}

export interface IAuthRegisterStudentParams {
  studentDraft: IStudentDraft
}

export interface IAuthValidateStudentLoginParams {
  studentLogin: IStudentLogin
}

export interface IAuthValidateStudentParams {
  studentDraft: IStudentDraft
}

export interface ICourse {
  courseId: number,
  title: string,
  description: string,
  timeLimit?: number,
  availableType: string,
  availableFrom?: string,
  availableTo?: string
}

export interface ICourseByIDParams {
  courseID: number
}

export interface ICourseListParams {
  page: number,
  pageSize: number
}

export interface ICourseSummary {
  courseId: number,
  title: string,
  timeLimit?: number,
  availableType: string,
  availableFrom?: string,
  availableTo?: string
}

export interface IExamAnswerParams {
  examID: number,
  questionID: number,
  optionIDs: Array<number>
}

export interface IExamGetQuestionParams {
  examID: number,
  questionID: number
}

export interface IExamHistoryParams {
  page: number,
  pageSize: number
}

export interface IExamResult {
  examId: number,
  status: string,
  finalScore: number,
  correctAnswers: number,
  totalQuestions: number
}

export interface IExamStart {
  examId: number,
  questionIds: Array<number>,
  startedAt: string,
  finishedAt?: string
}

export interface IExamStartParams {
  courseID: number
}

export interface IExamSubmitParams {
  examID: number
}

export interface IExamSummary {
  examId: number,
  courseId: number,
  status: string,
  finalScore: number,
  finishedAt: string
}

export interface IFieldError {
  field: string,
  error: string,
  constraint?: IFieldErrorConstraint
}

export interface IFieldErrorConstraint {
  max: number,
  min: number
}

export interface IQuestion {
  questionId: number,
  questionText: string,
  questionType: string,
  photoUrl?: string,
  options: Array<IQuestionOption>
}

export interface IQuestionOption {
  optionId: number,
  optionText: string
}

export interface IStudent {
  studentId: number,
  login: string,
  email: string,
  firstName: string,
  lastName: string
}

export interface IStudentDraft {
  login: string,
  email: string,
  password: string,
  firstName: string,
  lastName: string
}

export interface IStudentLogin {
  login: string,
  password: string
}

export interface IToken {
  accessToken: string,
  expiresIn: number,
  tokenType: string
}

export const factory = (send: any) => ({
  auth: {
    login(params: IAuthLoginParams): Promise<IToken> {
      return send('auth.Login', params)
    },
    registerStudent(params: IAuthRegisterStudentParams): Promise<IToken> {
      return send('auth.RegisterStudent', params)
    },
    validateStudent(params: IAuthValidateStudentParams): Promise<Array<IFieldError>> {
      return send('auth.ValidateStudent', params)
    },
    validateStudentLogin(params: IAuthValidateStudentLoginParams): Promise<Array<IFieldError>> {
      return send('auth.ValidateStudentLogin', params)
    }
  },
  course: {
    byID(params: ICourseByIDParams): Promise<ICourse> {
      return send('course.ByID', params)
    },
    list(params: ICourseListParams): Promise<Array<ICourseSummary>> {
      return send('course.List', params)
    },
    me(): Promise<IStudent> {
      return send('course.Me')
    }
  },
  exam: {
    answer(params: IExamAnswerParams): Promise<void> {
      return send('exam.Answer', params)
    },
    getQuestion(params: IExamGetQuestionParams): Promise<IQuestion> {
      return send('exam.GetQuestion', params)
    },
    history(params: IExamHistoryParams): Promise<Array<IExamSummary>> {
      return send('exam.History', params)
    },
    start(params: IExamStartParams): Promise<IExamStart> {
      return send('exam.Start', params)
    },
    submit(params: IExamSubmitParams): Promise<IExamResult> {
      return send('exam.Submit', params)
    }
  }
})

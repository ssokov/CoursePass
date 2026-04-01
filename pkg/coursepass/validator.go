package coursepass

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	FieldErrorRequired = "required"
	FieldErrorMax      = "max"
	FieldErrorMin      = "min"
	FieldErrorFormat   = "format"
	FieldErrorLen      = "len"

	fieldPathSeparator = "."
)

var errorMap = map[string]string{
	"required": FieldErrorRequired,
	"max":      FieldErrorMax,
	"min":      FieldErrorMin,
	"len":      FieldErrorLen,
	"email":    FieldErrorFormat,
}

var validate = newPlaygroundValidator()

type FieldError struct {
	Field      string                `json:"field"`
	Error      string                `json:"error"`
	Constraint *FieldErrorConstraint `json:"constraint,omitempty"`
}

type FieldErrorConstraint struct {
	Max int `json:"max,omitempty"`
	Min int `json:"min,omitempty"`
}

type FieldErrorConstraintFunc func(*FieldErrorConstraint)

type ValidationErrors []FieldError

func (e ValidationErrors) Error() string {
	return "validation error"
}

func (e ValidationErrors) Unwrap() error {
	return ErrValidation
}

type Validator struct {
	fields []FieldError
	err    error
}

func newPlaygroundValidator() *validator.Validate {
	vl := validator.New()
	vl.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name != "" {
			return name
		}
		return fld.Name
	})
	return vl
}

func NewFieldErrorConstraint(errorName, param string) *FieldErrorConstraint {
	value, err := strconv.Atoi(param)
	if err != nil {
		return nil
	}

	switch errorName {
	case FieldErrorMin:
		return &FieldErrorConstraint{Min: value}
	case FieldErrorMax:
		return &FieldErrorConstraint{Max: value}
	case FieldErrorLen:
		return &FieldErrorConstraint{Min: value, Max: value}
	default:
		return nil
	}
}

func NewFieldError(e validator.FieldError) FieldError {
	tag := e.Tag()
	if mappedError, ok := errorMap[tag]; ok {
		tag = mappedError
	}

	splitted := strings.Split(e.Namespace(), fieldPathSeparator)
	field := e.Field()
	if len(splitted) > 1 {
		field = strings.Join(splitted[1:], fieldPathSeparator)
	}

	return FieldError{
		Field:      field,
		Error:      tag,
		Constraint: NewFieldErrorConstraint(tag, e.Param()),
	}
}

func (v *Validator) SetInternalError(err error) {
	v.err = err
}

func (v *Validator) Append(field, errCode string, constraintFuncs ...FieldErrorConstraintFunc) {
	f := FieldError{Field: field, Error: errCode}
	if len(constraintFuncs) > 0 {
		c := &FieldErrorConstraint{}
		for _, fn := range constraintFuncs {
			fn(c)
		}
		f.Constraint = c
	}
	v.fields = append(v.fields, f)
}

func (v *Validator) CheckBasic(ctx context.Context, item any) {
	v.SetInternalError(nil)
	err := validate.StructCtx(ctx, item)
	if err == nil {
		return
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			v.fields = append(v.fields, NewFieldError(fieldError))
		}
		return
	}

	v.SetInternalError(err)
}

func (v *Validator) HasInternalError() bool {
	return v.err != nil
}

func (v *Validator) HasErrors() bool {
	return len(v.fields) != 0 || v.HasInternalError()
}

func (v *Validator) Fields() []FieldError {
	if len(v.fields) == 0 {
		return []FieldError{}
	}
	return v.fields
}

func (v *Validator) Error() error {
	if v.err != nil {
		return v.err
	}
	if len(v.fields) == 0 {
		return nil
	}
	return ValidationErrors(v.fields)
}

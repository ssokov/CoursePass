package rpc

import (
	"context"
	"errors"

	"courses/pkg/coursepass"
	"courses/pkg/coursepass/auth"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type AuthService struct {
	zenrpc.Service
	embedlog.Logger

	authManager *auth.Manager
}

func NewAuthService(dbc db.DB, logger embedlog.Logger, jwtSecret string, jwtTTLSeconds int) *AuthService {
	return &AuthService{
		authManager: auth.NewManager(dbc, logger, jwtSecret, jwtTTLSeconds),
		Logger:      logger,
	}
}

//zenrpc:500 internal error
func (as *AuthService) ValidateStudent(ctx context.Context, studentDraft StudentDraft) ([]FieldError, error) {
	err := as.authManager.ValidateStudent(ctx, studentDraft.ToModel())
	if err != nil {
		var validationErrs coursepass.ValidationErrors
		if errors.As(err, &validationErrs) {
			return newFieldErrors(validationErrs), nil
		}
		return nil, newInternalError(err)
	}

	return nil, nil
}

//zenrpc:409 login or email exists
//zenrpc:500 internal error
func (as *AuthService) RegisterStudent(ctx context.Context, studentDraft StudentDraft) (*Token, error) {
	token, err := as.authManager.RegisterStudent(ctx, studentDraft.ToModel())
	if err != nil {
		if errors.Is(err, coursepass.ErrValidation) {
			return nil, ErrInvalidParams
		}
		if errors.Is(err, coursepass.ErrLoginExists) {
			return nil, ErrLoginExists
		}
		if errors.Is(err, coursepass.ErrEmailExists) {
			return nil, ErrEmailExists
		}
		return nil, newInternalError(err)
	}

	return newToken(token), nil
}

//zenrpc:500 internal error
func (as *AuthService) ValidateStudentLogin(ctx context.Context, studentLogin StudentLogin) ([]FieldError, error) {
	err := as.authManager.ValidateStudentLogin(ctx, studentLogin.ToModel())
	if err != nil {
		var validationErrs coursepass.ValidationErrors
		if errors.As(err, &validationErrs) {
			return newFieldErrors(validationErrs), nil
		}
		return nil, newInternalError(err)
	}

	return nil, nil
}

//zenrpc:401 invalid credentials
//zenrpc:500 internal error
func (as *AuthService) Login(ctx context.Context, studentLogin StudentLogin) (*Token, error) {
	token, err := as.authManager.Login(ctx, studentLogin.ToModel())
	if err != nil {
		var validationErrs coursepass.ValidationErrors
		if errors.As(err, &validationErrs) && len(validationErrs) > 0 {
			return nil, newInvalidParamsError(validationErrs[0].Field, validationErrs[0].Error)
		}
		if errors.Is(err, coursepass.ErrInvalidCredentials) {
			return nil, ErrInvalidCredentials
		}
		return nil, newInternalError(err)
	}

	return newToken(token), nil
}

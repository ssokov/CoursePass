package rpc

import (
	"context"
	"net/mail"

	"courses/pkg/coursepass"
	"courses/pkg/db"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type AuthService struct {
	zenrpc.Service
	embedlog.Logger

	authManager *coursepass.AuthManager
}

func NewAuthService(dbc db.DB, logger embedlog.Logger, authCfg coursepass.AuthConfig) *AuthService {
	return &AuthService{
		authManager: coursepass.NewAuthManager(dbc, logger, authCfg),
		Logger:      logger,
	}
}

func (as *AuthService) Register(ctx context.Context, login, password, email, firstName, lastName string) (*Token, error) {
	if err := as.validateRegisterRequest(login, password, email, firstName, lastName); err != nil {
		return nil, err
	}

	token, err := as.authManager.Register(ctx, login, password, email, firstName, lastName)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newToken(token), nil
}

func (as *AuthService) Login(ctx context.Context, login, password string) (*Token, error) {
	if err := as.validateLoginRequest(login, password); err != nil {
		return nil, err
	}

	token, err := as.authManager.Login(ctx, login, password)
	if err != nil {
		return nil, newInternalError(err)
	}

	return newToken(token), nil
}

func (as *AuthService) validateRegisterRequest(login, password, email, firstName, lastName string) error {
	if login == "" {
		return newInvalidParamsError("login", "is required")
	}
	if len([]rune(login)) > 255 {
		return newInvalidParamsError("login", "max length is 255")
	}

	if password == "" {
		return newInvalidParamsError("password", "is required")
	}
	if len([]rune(password)) < 6 {
		return newInvalidParamsError("password", "min length is 6")
	}
	if len([]rune(password)) > 255 {
		return newInvalidParamsError("password", "max length is 255")
	}

	if email == "" {
		return newInvalidParamsError("email", "is required")
	}
	if len([]rune(email)) > 255 {
		return newInvalidParamsError("email", "max length is 255")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return newInvalidParamsError("email", "invalid format")
	}

	if firstName == "" {
		return newInvalidParamsError("firstName", "is required")
	}
	if len([]rune(firstName)) > 255 {
		return newInvalidParamsError("firstName", "max length is 255")
	}

	if lastName == "" {
		return newInvalidParamsError("lastName", "is required")
	}
	if len([]rune(lastName)) > 255 {
		return newInvalidParamsError("lastName", "max length is 255")
	}

	return nil
}

func (as *AuthService) validateLoginRequest(login, password string) error {
	if login == "" {
		return newInvalidParamsError("login", "is required")
	}
	if len([]rune(login)) > 255 {
		return newInvalidParamsError("login", "max length is 255")
	}

	if password == "" {
		return newInvalidParamsError("password", "is required")
	}
	if len([]rune(password)) > 255 {
		return newInvalidParamsError("password", "max length is 255")
	}

	return nil
}

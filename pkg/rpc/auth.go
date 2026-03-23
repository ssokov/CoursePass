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

func (as *AuthService) Register(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		as.Logger.Error(ctx, "auth register invalid params", "err", err)
		return RegisterResponse{}, err
	}

	token, err := as.authManager.Register(
		ctx,
		req.Login,
		req.Password,
		req.Email,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		as.Logger.Error(ctx, "auth register failed", "err", err)
		return RegisterResponse{}, mapRPCError(err)
	}

	return newRegisterResponse(token), nil
}

func (as *AuthService) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	if err := validateLoginRequest(req); err != nil {
		as.Logger.Error(ctx, "auth login invalid params", "err", err)
		return LoginResponse{}, err
	}

	token, err := as.authManager.Login(ctx, req.Login, req.Password)
	if err != nil {
		as.Logger.Error(ctx, "auth login failed", "err", err)
		return LoginResponse{}, mapRPCError(err)
	}

	return newLoginResponse(token), nil
}

func validateRegisterRequest(req RegisterRequest) error {
	if req.Login == "" {
		return invalidParamsError("login", "is required")
	}
	if len([]rune(req.Login)) > 255 {
		return invalidParamsError("login", "max length is 255")
	}

	if req.Password == "" {
		return invalidParamsError("password", "is required")
	}
	if len([]rune(req.Password)) < 6 {
		return invalidParamsError("password", "min length is 6")
	}
	if len([]rune(req.Password)) > 255 {
		return invalidParamsError("password", "max length is 255")
	}

	if req.Email == "" {
		return invalidParamsError("email", "is required")
	}
	if len([]rune(req.Email)) > 255 {
		return invalidParamsError("email", "max length is 255")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return invalidParamsError("email", "invalid format")
	}

	if req.FirstName == "" {
		return invalidParamsError("firstName", "is required")
	}
	if len([]rune(req.FirstName)) > 255 {
		return invalidParamsError("firstName", "max length is 255")
	}

	if req.LastName == "" {
		return invalidParamsError("lastName", "is required")
	}
	if len([]rune(req.LastName)) > 255 {
		return invalidParamsError("lastName", "max length is 255")
	}

	return nil
}

func validateLoginRequest(req LoginRequest) error {
	if req.Login == "" {
		return invalidParamsError("login", "is required")
	}
	if len([]rune(req.Login)) > 255 {
		return invalidParamsError("login", "max length is 255")
	}

	if req.Password == "" {
		return invalidParamsError("password", "is required")
	}
	if len([]rune(req.Password)) > 255 {
		return invalidParamsError("password", "max length is 255")
	}

	return nil
}

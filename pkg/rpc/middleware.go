package rpc

import (
	"context"
	"encoding/json"
	"strings"

	"courses/pkg/coursepass"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type studentCtx string

const (
	studentKey   studentCtx = "rpc.student.id"
	bearerPrefix string     = "Bearer "
)

func authMiddleware(authCfg coursepass.AuthConfig, logger embedlog.Logger) zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			_, ok := zenrpc.RequestFromContext(ctx)
			if !ok {
				return h(ctx, method, params)
			}

			ns := zenrpc.NamespaceFromContext(ctx)
			if ns == NSAuth {
				return h(ctx, method, params)
			}

			token, err := bearerTokenFromContext(ctx)
			if err != nil {
				logger.Error(ctx, "auth middleware: invalid token", "err", err)
				return zenrpc.NewResponseError(zenrpc.IDFromContext(ctx), errInvalidToken, "invalid token", nil)
			}

			studentID, err := coursepass.ValidateJWT(authCfg, token)
			if err != nil {
				logger.Error(ctx, "auth middleware: token validation failed", "err", err)
				return zenrpc.NewResponseError(zenrpc.IDFromContext(ctx), errInvalidToken, "invalid token", nil)
			}

			return h(context.WithValue(ctx, studentKey, studentID), method, params)
		}
	}
}

func StudentIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(studentKey).(int)
	return id, ok
}

func bearerTokenFromContext(ctx context.Context) (string, error) {
	req, ok := zenrpc.RequestFromContext(ctx)
	if !ok {
		return "", coursepass.ErrInvalidToken
	}

	auth := req.Header.Get("Authorization")
	if auth == "" {
		return "", coursepass.ErrInvalidToken
	}

	if !strings.HasPrefix(auth, bearerPrefix) {
		return "", coursepass.ErrInvalidToken
	}

	token := strings.TrimSpace(strings.TrimPrefix(auth, bearerPrefix))
	if token == "" {
		return "", coursepass.ErrInvalidToken
	}

	return token, nil
}

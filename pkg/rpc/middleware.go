package rpc

import (
	"context"
	"encoding/json"
	"strings"

	"courses/pkg/course"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type studentCtx string

const (
	studentKey   studentCtx = "rpc.student.id"
	bearerPrefix string     = "Bearer "
)

func authMiddleware(authCfg course.AuthConfig, logger embedlog.Logger) zenrpc.MiddlewareFunc {
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

			studentID, err := course.ValidateJWT(authCfg, token)
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
		return "", course.ErrInvalidToken
	}

	auth := req.Header.Get("Authorization")
	if auth == "" {
		return "", course.ErrInvalidToken
	}

	if !strings.HasPrefix(auth, bearerPrefix) {
		return "", course.ErrInvalidToken
	}

	token := strings.TrimSpace(strings.TrimPrefix(auth, bearerPrefix))
	if token == "" {
		return "", course.ErrInvalidToken
	}

	return token, nil
}

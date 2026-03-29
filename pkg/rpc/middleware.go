package rpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"courses/pkg/coursepass"

	"github.com/vmkteam/embedlog"
	"github.com/vmkteam/zenrpc/v2"
)

type studentCtx string

const (
	studentKey   studentCtx = "rpc.student.id"
	bearerPrefix string     = "Bearer "
)

type tokenClaims struct {
	Sub   string
	Login string
	Exp   int64
	Iat   int64
}

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
				return zenrpc.NewResponseError(
					zenrpc.IDFromContext(ctx),
					ErrInvalidToken.Code,
					ErrInvalidToken.Message,
					ErrInvalidToken.Data,
				)
			}

			studentID, err := validateJWT(authCfg.JWTSecret, token)
			if err != nil {
				logger.Error(ctx, "auth middleware: token validation failed", "err", err)
				return zenrpc.NewResponseError(
					zenrpc.IDFromContext(ctx),
					ErrInvalidToken.Code,
					ErrInvalidToken.Message,
					ErrInvalidToken.Data,
				)
			}

			return h(context.WithValue(ctx, studentKey, studentID), method, params)
		}
	}
}

func studentIDFromContext(ctx context.Context) (int, bool) {
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

func validateJWT(jwtSecretValue, token string) (int, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSig := signHS256(unsigned, jwtSecret(jwtSecretValue))
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return 0, ErrInvalidToken
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, ErrInvalidToken
	}

	claims, err := unmarshalTokenClaims(payloadRaw)
	if err != nil {
		return 0, ErrInvalidToken
	}

	if claims.Sub == "" || claims.Exp <= 0 {
		return 0, ErrInvalidToken
	}
	if time.Now().Unix() >= claims.Exp {
		return 0, ErrInvalidToken
	}

	studentID, err := strconv.Atoi(claims.Sub)
	if err != nil || studentID <= 0 {
		return 0, ErrInvalidToken
	}

	return studentID, nil
}

func signHS256(unsignedToken string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsignedToken))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func unmarshalTokenClaims(data []byte) (tokenClaims, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return tokenClaims{}, err
	}

	sub, ok := raw["sub"].(string)
	if !ok || sub == "" {
		return tokenClaims{}, ErrInvalidToken
	}

	login, _ := raw["login"].(string)

	exp, err := numberFieldAsInt64(raw["exp"])
	if err != nil {
		return tokenClaims{}, err
	}

	iat, err := numberFieldAsInt64(raw["iat"])
	if err != nil {
		return tokenClaims{}, err
	}

	return tokenClaims{
		Sub:   sub,
		Login: login,
		Exp:   exp,
		Iat:   iat,
	}, nil
}

func jwtSecret(secret string) []byte {
	if secret == "" {
		secret = "coursepass-dev-secret"
	}
	return []byte(secret)
}

func numberFieldAsInt64(v any) (int64, error) {
	switch n := v.(type) {
	case float64:
		return int64(n), nil
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	default:
		return 0, ErrInvalidToken
	}
}

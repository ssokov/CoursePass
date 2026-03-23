package coursepass

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func generateJWT(authCfg AuthConfig, studentID int, login string) (string, int, error) {
	header := newTokenHeader()

	now := time.Now()
	ttl := tokenTTLSeconds(authCfg)
	exp := now.Unix() + int64(ttl)
	claims := newTokenClaims(studentID, login, now.Unix(), exp)

	headerRaw, err := marshalTokenHeader(header)
	if err != nil {
		return "", 0, err
	}
	claimsRaw, err := marshalTokenClaims(claims)
	if err != nil {
		return "", 0, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerRaw)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsRaw)
	unsigned := encodedHeader + "." + encodedClaims
	signature := signHS256(unsigned, jwtSecret(authCfg))

	return unsigned + "." + signature, ttl, nil
}

func ValidateJWT(authCfg AuthConfig, token string) (int, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	expectedSig := signHS256(unsigned, jwtSecret(authCfg))
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

func marshalTokenHeader(h tokenHeader) ([]byte, error) {
	return json.Marshal(map[string]any{
		"alg": h.Alg,
		"typ": h.Typ,
	})
}

func marshalTokenClaims(c tokenClaims) ([]byte, error) {
	return json.Marshal(map[string]any{
		"sub":   c.Sub,
		"login": c.Login,
		"exp":   c.Exp,
		"iat":   c.Iat,
	})
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

func jwtSecret(authCfg AuthConfig) []byte {
	secret := authCfg.JWTSecret
	if secret == "" {
		secret = "coursepass-dev-secret"
	}
	return []byte(secret)
}

func tokenTTLSeconds(authCfg AuthConfig) int {
	if authCfg.JWTTTLSeconds <= 0 {
		return defaultTokenTTLSeconds
	}

	return authCfg.JWTTTLSeconds
}

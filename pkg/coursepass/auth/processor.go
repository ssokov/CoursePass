package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	"courses/pkg/coursepass"
)

type processor struct {
	secretValue string
	ttlSeconds  int
}

func newAuthProcessor(secretValue string, ttlSeconds int) *processor {
	return &processor{secretValue: secretValue, ttlSeconds: ttlSeconds}
}

func (p *processor) generateToken(studentID int, login string) (*coursepass.AuthToken, error) {
	token, expiresIn, err := p.generateJWT(studentID, login)
	if err != nil {
		return nil, err
	}

	return coursepass.NewAuthToken(token, expiresIn), nil
}

func (p *processor) generateJWT(studentID int, login string) (string, int, error) {
	header := coursepass.NewTokenHeader()

	now := time.Now()
	ttl := p.tokenTTLSeconds()
	exp := now.Unix() + int64(ttl)
	claims := coursepass.NewTokenClaims(studentID, login, now.Unix(), exp)

	headerRaw, err := p.marshalTokenHeader(header)
	if err != nil {
		return "", 0, err
	}
	claimsRaw, err := p.marshalTokenClaims(claims)
	if err != nil {
		return "", 0, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerRaw)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsRaw)
	unsigned := encodedHeader + "." + encodedClaims
	signature := p.signHS256(unsigned, p.jwtSecret())

	return unsigned + "." + signature, ttl, nil
}

func (p *processor) signHS256(unsignedToken string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(unsignedToken))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (p *processor) marshalTokenClaims(c coursepass.TokenClaims) ([]byte, error) {
	return json.Marshal(map[string]any{
		"sub":   c.Sub,
		"login": c.Login,
		"exp":   c.Exp,
		"iat":   c.Iat,
	})
}

func (p *processor) marshalTokenHeader(h coursepass.TokenHeader) ([]byte, error) {
	return json.Marshal(map[string]any{
		"alg": h.Alg,
		"typ": h.Typ,
	})
}

func (p *processor) jwtSecret() []byte {
	if p.secretValue == "" {
		return []byte("coursepass-dev-secret")
	}
	return []byte(p.secretValue)
}

func (p *processor) tokenTTLSeconds() int {
	if p.ttlSeconds <= 0 {
		return coursepass.DefaultTokenTTLSeconds
	}

	return p.ttlSeconds
}

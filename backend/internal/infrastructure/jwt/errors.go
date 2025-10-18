package jwt

import "github.com/cockroachdb/errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidSecret    = errors.New("invalid secret")
	ErrInvalidClaims    = errors.New("invalid claims")
)

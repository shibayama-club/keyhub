package console

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/shibayama-club/keyhub/internal/domain/authenticator"
	"github.com/shibayama-club/keyhub/internal/infrastructure/auth/claim"
	"github.com/shibayama-club/keyhub/internal/infrastructure/jwt"
)

type AuthService struct {
	generator *jwt.Generator
	validator *jwt.Validator
}

var _ authenticator.ConsoleAuthenticator = (*AuthService)(nil)

func NewAuthService(secret string) (*AuthService, error) {
	generator, err := jwt.NewGenerator(secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create JWT generator")
	}

	validator, err := jwt.NewValidator(secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create JWT validator")
	}

	return &AuthService{
		generator: generator,
		validator: validator,
	}, nil
}

func (s *AuthService) GenerateToken(organizationID, sessionID string, expiresIn time.Duration) (string, error) {
	claims := claim.NewConsoleClaims(organizationID, sessionID)
	return s.generator.Generate(claims, expiresIn)
}

func (s *AuthService) ValidateToken(token string) (*claim.ConsoleClaims, error) {
	claims := &claim.ConsoleClaims{}
	if err := s.validator.Validate(token, claims); err != nil {
		return nil, err
	}
	return claims, nil
}

package authenticator

import (
	"time"

	"github.com/shibayama-club/keyhub/internal/infrastructure/auth/claim"
)

type ConsoleAuthenticator interface {
	GenerateToken(organizationID, sessionID string, expiresIn time.Duration) (string, error)
	ValidateToken(token string) (*claim.ConsoleClaims, error)
}

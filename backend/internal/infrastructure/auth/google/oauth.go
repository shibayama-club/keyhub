package google

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"

	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
	"google.golang.org/api/idtoken"
)

type OAuthService struct {
	clientID     string
	clientSecret string
	redirectURI  string
	httpClient   *http.Client
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

const (
	GoogleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleJWKSURL  = "https://www.googleapis.com/oauth2/v3/certs"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type IDTokenClaims struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
	Nonce         string `json:"nonce"`
}

func NewOAuthService(config OAuthConfig) (*OAuthService, error) {
	if config.ClientID == "" {
		return nil, errors.New("client ID is required")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("client secret is required")
	}
	if config.RedirectURI == "" {
		return nil, errors.New("redirect URI is required")
	}

	return &OAuthService{
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
		redirectURI:  config.RedirectURI,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func GenerateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "failed to generate code verifier")
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "failed to generate random string")
	}
	return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

func (s *OAuthService) BuildAuthURL(state, nonce, codeChallenge string) string {
	params := url.Values{}
	params.Set("client_id", s.clientID)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	params.Set("nonce", nonce)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	return GoogleAuthURL + "?" + params.Encode()
}

func (s *OAuthService) ExchangeCode(ctx context.Context, code, codeVerifier string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", s.clientID)
	data.Set("client_secret", s.clientSecret)
	data.Set("redirect_uri", s.redirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequestWithContext(ctx, "POST", GoogleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create token request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange code for token")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Newf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, errors.Wrap(err, "failed to decode token response")
	}

	return &tokenResp, nil
}

func (s *OAuthService) VerifyIDToken(ctx context.Context, idTokenString, expectedNonce string) (*IDTokenClaims, error) {
	payload, err := idtoken.Validate(ctx, idTokenString, s.clientID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate ID token")
	}

	claims := IDTokenClaims{
		Iss:           payload.Issuer,
		Aud:           payload.Audience,
		Sub:           payload.Subject,
		Exp:           payload.Expires,
		Iat:           payload.IssuedAt,
		Email:         getStringClaim(payload.Claims, "email"),
		EmailVerified: getBoolClaim(payload.Claims, "email_verified"),
		Name:          getStringClaim(payload.Claims, "name"),
		Picture:       getStringClaim(payload.Claims, "picture"),
		GivenName:     getStringClaim(payload.Claims, "given_name"),
		FamilyName:    getStringClaim(payload.Claims, "family_name"),
		Locale:        getStringClaim(payload.Claims, "locale"),
		Nonce:         getStringClaim(payload.Claims, "nonce"),
		Azp:           getStringClaim(payload.Claims, "azp"),
	}

	if claims.Nonce != expectedNonce {
		return nil, errors.New("invalid nonce")
	}

	if !claims.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return &claims, nil
}

func getStringClaim(claims map[string]any, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}

func getBoolClaim(claims map[string]any, key string) bool {
	if val, ok := claims[key].(bool); ok {
		return val
	}
	return false
}

func (c *IDTokenClaims) GetUserInfo() (email, name, picture, sub string) {
	return c.Email, c.Name, c.Picture, c.Sub
}

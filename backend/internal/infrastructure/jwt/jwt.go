package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Claims interface {
	// 有効期限を取得+トークン検証時に期限切れをチェック
	GetExpiration() int64
	// 発行時刻を取得+トークンの発行時刻を確認
	GetIssuedAt() int64
	// 有効期限を設定+トークン生成時に自動設定
	SetExpiration(exp int64)
	// 発行時刻を設定+トークン生成時に自動設定
	SetIssuedAt(iat int64)
}

type Generator struct {
	secret string
}

func NewGenerator(secret string) (*Generator, error) {
	if secret == "" {
		return nil, ErrInvalidSecret
	}
	return &Generator{
		secret: secret,
	}, nil
}

func (g *Generator) Generate(claims Claims, expiresIn time.Duration) (string, error) {
	if claims == nil {
		return "", ErrInvalidClaims
	}

	if expiresIn <= 0 {
		return "", errors.Wrap(ErrInvalidClaims, "expiration duration must be positive")
	}

	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	now := time.Now()
	claims.SetIssuedAt(now.Unix())
	claims.SetExpiration(now.Add(expiresIn).Unix())

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal header")
	}
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal claims")
	}
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	message := headerB64 + "." + claimsB64
	h := hmac.New(sha256.New, []byte(g.secret))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	token := message + "." + signature

	return token, nil
}

type Validator struct {
	secret string
}

func NewValidator(secret string) (*Validator, error) {
	if secret == "" {
		return nil, ErrInvalidSecret
	}
	return &Validator{
		secret: secret,
	}, nil
}

func (v *Validator) Validate(tokenString string, claims Claims) error {
	if tokenString == "" {
		return ErrInvalidToken
	}

	if claims == nil {
		return ErrInvalidClaims
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return ErrInvalidToken
	}

	for _, part := range parts {
		if part == "" {
			return ErrInvalidToken
		}
	}

	message := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, []byte(v.secret))
	h.Write([]byte(message))
	expectedSig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return ErrInvalidSignature
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return errors.Wrap(ErrInvalidToken, "failed to decode claims")
	}

	if err := json.Unmarshal(claimsJSON, claims); err != nil {
		return errors.Wrap(ErrInvalidClaims, "failed to unmarshal claims")
	}

	if time.Now().Unix() > claims.GetExpiration() {
		return ErrTokenExpired
	}

	return nil
}

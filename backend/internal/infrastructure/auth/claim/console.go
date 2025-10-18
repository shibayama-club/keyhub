package claim

type ConsoleClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Org string `json:"org"`
	Sid string `json:"sid"`
}

func (c *ConsoleClaims) GetExpiration() int64 {
	return c.Exp
}

func (c *ConsoleClaims) GetIssuedAt() int64 {
	return c.Iat
}

func (c *ConsoleClaims) SetExpiration(exp int64) {
	c.Exp = exp
}

func (c *ConsoleClaims) SetIssuedAt(iat int64) {
	c.Iat = iat
}

func NewConsoleClaims(organizationID, sessionID string) *ConsoleClaims {
	return &ConsoleClaims{
		Sub: organizationID,
		Org: organizationID,
		Sid: sessionID,
	}
}

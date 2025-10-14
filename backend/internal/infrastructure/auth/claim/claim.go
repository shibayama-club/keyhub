package claim

type Claims interface {
	GetExpiration() int64
	GetIssuedAt() int64
	SetExpiration(exp int64)
	SetIssuedAt(iat int64)
}

package model

import (
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type TenantID uuid.UUID

func (id TenantID) UUID() uuid.UUID{
	return uuid.UUID(id)
}

func (id TenantID) String() string{
	return uuid.UUID(id).String()
}

type TenantName string

func (n TenantName) String() string{
	return string(n)
}

func (n TenantName) Validate() error{
	if n == ""{
		return errors.WithHint(
			errors.New("tenant name is required"),
			"Tenant Nameは必須です。",
		)
	}

	if utf8.RuneCountInString(string(n)) > 30{
		return errors.WithHint(
			errors.New("Please enter a tenantname within 30 characters"),
			"テナント名は30文字以内で入力してください。",
		)
	}
	return nil
}

func NewTenantName(value string) (TenantName, error){
	n := TenantName(value)
	if err := n.Validate(); err != nil{
		return "", err
	}
	return n, nil
}

type TenantSlug string

func (s TenantSlug) String() string{
	return string(s)
}

func (s TenantSlug) Validate() error{
	if s == ""{
		return nil
	}
	if utf8.RuneCountInString(string(s)) > 30 {
        return errors.WithHint(
            errors.New("please enter a tenant slug within 30 characters"),
            "Slugは30文字以内で入力してください。",
        )
    }
	if !IsSlugFormat(string(s)){
		return errors.WithHint(
			errors.New("Please enter a validate slug address"),
			"Slugの正しい形式で入力してください",
		)
	}
	return nil
}

func NewTenantSlug(value string)(TenantSlug, error){
	s := TenantSlug(value)
	if err := s.Validate(); err != nil{
		return "", err
	}
	return s, nil
}

type TenantPasswordHash string

func(p TenantPasswordHash) String() string{
	return string(p)
}

func(p TenantPasswordHash) Validate() error{
	if p == ""{
		return errors.WithHint(
			errors.New("tenant password hash is required"),
			"テナントパスワードハッシュは必須です。",
		)
	}
	return nil
}

func NewTenantPasswordHash(value string)(TenantPasswordHash, error){
	p := TenantPasswordHash(value)
	if err := p.Validate(); err != nil{
		return "", err
	}
	return p, nil
}

type Tenant struct{
	TenantId TenantID
	TenantName TenantName
	TenantSlug TenantSlug
	TenantPasswordHash TenantPasswordHash
	CreatedAt time.Time
	UpdatedAt time.Time
}

var isSlugFormatRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func IsSlugFormat(str string) bool{
	return isSlugFormatRegex.MatchString(str)
}



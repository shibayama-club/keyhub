package console

import (
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type ConsoleID uuid.UUID

func (id ConsoleID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id ConsoleID) String() string {
	return uuid.UUID(id).String()
}

type ConsoleEmail string

func (e ConsoleEmail) String() string {
	return string(e)
}

func (e ConsoleEmail) Validate() error {
	if e =="" {
		return errors.WithHint(
			errors.New("email is required"),
			"Emailは必須です。",
		)
	}
	
	if !IsEmailFormat(string(e)){
		return errors.WithHint(
			errors.New("Please enter a validate email address"),
			"Emailの正しい形式で入力してください。",
		)
	}
	return nil
}

func NewConsoleEmail(value string)(ConsoleEmail, error){
	e := ConsoleEmail(value)
	if err := e.Validate(); err != nil{
		return "", err
	}
	return e, nil
}

type ConsoleName string

func (n ConsoleName) String() string {
	return string(n)
}

func (n ConsoleName) Validate() error{
	if n == "" {
		return errors.WithHint(
			errors.New("consolename is required"),
			"コンソールネームは必須です。",
		)
	}
	if utf8.RuneCountInString(string(n)) > 30 {
		return errors.WithHint(
			errors.New("Please enter a console name within 10 characters"),
			"コンソールネームは10字以内で入力してください。",
		)
	}
	return nil
}

func NewConsoleName(value string)(ConsoleName, error){
	n := ConsoleName(value)
	if err := n.Validate(); err != nil{
		return "", err
	}
	return n, nil
}	

type ConsoleIcon string

func (i ConsoleIcon) Validate() error{
	if i == ""{
		return nil
	}

	if !IsIconURLFormat(string(i)){
		return errors.WithHint(
			errors.New("icon must be a valid HTTPS URL"),
			"アイコンURLは有効なHTTPS形式である必要があります。",
		)
	}

	if utf8.RuneCountInString(string(i)) > 2048 {
		return errors.WithHint(
			errors.New("icon URL is too long"),
			"アイコンURLが長すぎます。",
		)
	}

	return nil
}

func NewConsoleIcon(value string)(ConsoleIcon, error){
	i := ConsoleIcon(value)
	if err := i.Validate(); err != nil{
		return "", err
	}
	return i, nil
}

type CampusID uuid.UUID

func (id CampusID) UUID() uuid.UUID {
	return uuid.UUID(id)
}

func (id CampusID) String() string {
	return uuid.UUID(id).String()
}

type Console struct{
    Id ConsoleID
    Email ConsoleEmail
    Name ConsoleName
    Icon ConsoleIcon
    CampusId CampusID
    CreatedAt time.Time
    UpdatedAt time.Time
}

var isEmailFormatRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func IsEmailFormat(str string) bool {
	return isEmailFormatRegex.MatchString(str)
}

var isIconURLFormatRegex = regexp.MustCompile(`^https://[^\s]+$`)

func IsIconURLFormat(str string) bool {
	return isIconURLFormatRegex.MatchString(str)
}

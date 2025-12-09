package errors

import "github.com/cockroachdb/errors"

var (
	ErrValidation    = errors.New("Validation Error")
	ErrNotFound      = errors.New("Not Found Error")
	ErrUnAuthorized  = errors.New("Unauthorized Error")
	ErrInternal      = errors.New("Internal Error")
	ErrAlreadyExists = errors.New("Already Exists Error")
)

func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsUnAuthorizedError(err error) bool {
	return errors.Is(err, ErrUnAuthorized)
}

func IsInternalError(err error) bool {
	return errors.Is(err, ErrInternal)
}

func IsAlreadyExistsError(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

package errors

import "github.com/cockroachdb/errors"

var (
	ErrValidation = errors.New("Validation error")
	ErrNotFound = errors.New("not found error")
	ErrUnAuthorized = errors.New("unauthorized error")
	ErrInternal = errors.New("internal error")
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
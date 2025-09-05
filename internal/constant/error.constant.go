package constant

import "errors"

var (
	ErrNotFound     = errors.New("requested resource not found")
	ErrUnauthorized = errors.New("invalid credentials or token")
	ErrForbidden    = errors.New("user does not have the required permission")
	ErrConflict     = errors.New("resource already exists")
	ErrInvalidInput = errors.New("invalid input provided")
)

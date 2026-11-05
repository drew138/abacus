package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
	ErrInvalid     = errors.New("invalid")
	ErrUnavailable = errors.New("unavailable")
	ErrInternal    = errors.New("internal")
)

var (
	ErrUnknownCategory    = fmt.Errorf("unknown category: %w", ErrInvalid)
	ErrInvalidCredentials = fmt.Errorf("invalid credentials: %w", ErrInvalid)
	ErrNotAuthenticated   = fmt.Errorf("not authenticated: %w", ErrInvalid)
	ErrPasswordTooShort   = fmt.Errorf("password too short: %w", ErrInvalid)
)

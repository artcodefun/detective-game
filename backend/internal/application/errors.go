package application

import (
	"errors"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
)

type ErrorKind int

const (
	KindInternal ErrorKind = iota
	KindNotFound
	KindForbidden
	KindConflict
	KindInvalidInput
)

type AppError struct {
	Kind   ErrorKind
	Code   string
	Params map[string]any
}

func (e AppError) Error() string {
	return e.Code
}

func NewAppError(kind ErrorKind, code string) AppError {
	return AppError{Kind: kind, Code: code}
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	var appErr AppError
	if errors.As(err, &appErr) {
		return err
	}
	if errors.Is(err, ports.ErrNotFound) {
		return ErrNotFound
	}
	return NewAppError(KindInternal, "internal_error")
}

var (
	ErrNotFound     = NewAppError(KindNotFound, "not_found")
	ErrForbidden    = NewAppError(KindForbidden, "forbidden")
	ErrConflict     = NewAppError(KindConflict, "conflict")
	ErrInvalidInput = NewAppError(KindInvalidInput, "invalid_input")
)

func IsNotFound(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && appErr.Kind == KindNotFound
}

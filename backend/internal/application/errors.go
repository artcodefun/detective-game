package application

import "errors"

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

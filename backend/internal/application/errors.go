package application

import (
	"errors"

	"github.com/artcodefun/detective-game/backend/internal/application/ports"
	"github.com/artcodefun/detective-game/backend/internal/domain"
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
	Kind        ErrorKind
	Translation domain.Translation
}

func (e AppError) Error() string {
	return string(e.Translation.Key)
}

func NewAppError(kind ErrorKind, translation domain.Translation) AppError {
	return AppError{Kind: kind, Translation: translation}
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
	return NewAppError(KindInternal, domain.T("error.internal"))
}

var (
	ErrNotFound     = NewAppError(KindNotFound, domain.T("error.not_found"))
	ErrForbidden    = NewAppError(KindForbidden, domain.T("error.forbidden"))
	ErrConflict     = NewAppError(KindConflict, domain.T("error.conflict"))
	ErrInvalidInput = NewAppError(KindInvalidInput, domain.T("error.invalid_input"))
)

func IsNotFound(err error) bool {
	var appErr AppError
	return errors.As(err, &appErr) && appErr.Kind == KindNotFound
}

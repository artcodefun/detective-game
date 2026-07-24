package application

import "github.com/google/uuid"

type Actor struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}

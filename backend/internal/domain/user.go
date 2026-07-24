package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser() User {
	return User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func NewUser() User {
	return User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}

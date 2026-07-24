package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func NewUser() User {
	return User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}
}

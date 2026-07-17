package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Password_hash string    `json:"password_hash"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Score         int       `json:"score"`
	Created_at    time.Time `json:"created_at"`
}

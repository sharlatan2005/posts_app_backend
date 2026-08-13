package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID            uuid.UUID   `json:"id"`
	Username      string      `json:"username"`
	Password_hash string      `json:"password_hash"`
	Name          pgtype.Text `json:"name"`
	Surname       pgtype.Text `json:"surname"`
	Score         int         `json:"score"`
	Created_at    time.Time   `json:"created_at"`
}

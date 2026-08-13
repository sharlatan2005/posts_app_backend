package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	ID        uuid.UUID   `json:"id"`
	AuthorID  uuid.UUID   `json:"author_id"`
	Text      string      `json:"text"`
	CreatedAt pgtype.Text `json:"created_at"`
}

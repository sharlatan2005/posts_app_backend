package domain

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Comment struct {
	ID        uuid.UUID   `json:"id"`
	PostID    uuid.UUID   `json:"post_id"`
	AuthorID  uuid.UUID   `json:"author_id"`
	Text      string      `json:"text"`
	CreatedAt pgtype.Text `json:"created_at"`
}

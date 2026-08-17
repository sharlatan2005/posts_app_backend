package domain

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID      uuid.UUID `json:"id"`
	PostID  uuid.UUID `json:"post_id"`
	LikerID uuid.UUID `json:"liker_id"`
	LikedAt time.Time `json:"liked_at"`
}

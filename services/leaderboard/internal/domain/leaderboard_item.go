package domain

type LeaderboardItem struct {
	UserID string `json:"user_id"`
	Score  int64  `json:"score"`
}

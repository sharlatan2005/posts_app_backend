package servErrors

import "errors"

var (
	ErrCommentTextEmpty = errors.New("Текст комментария не может быть пустым.")
)

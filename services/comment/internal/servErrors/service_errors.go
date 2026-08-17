package servErrors

import "errors"

var (
	ErrCommentTextEmpty       = errors.New("Текст комментария не может быть пустым.")
	ErrCommentNotFound        = errors.New("Комментария не существует.")
	ErrCommentForbiddenAction = errors.New("У вас нет прав на действие с этим комментарием.")

	ErrPostNotFound     = errors.New("Не существует такого поста.")
	ErrNoCommentsOnPost = errors.New("На посту нет комментариев или такого поста нет.")
)

package servErrors

import "errors"

var (
	ErrPostTextEmpty = errors.New("Текст поста не может быть пустым.")
	ErrPostNotFound  = errors.New("Пост не найден")
)

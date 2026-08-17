package servErrors

import "errors"

var (
	ErrLikeNotFound = errors.New("Пользователь не ставил лайков под этим постом.")

	ErrLikeForbiddenAction = errors.New("Вы не можете убрать лайк этого пользователя с поста.")

	ErrNoLikesOnPost = errors.New("На данном посту нет лайков или поста не существует.")
)

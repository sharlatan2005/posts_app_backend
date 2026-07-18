package servErrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("Пользователь с таким именем уже существует")
	ErrUserNotFound      = errors.New("Нет такого пользователя")
	ErrUsernameEmpty     = errors.New("Username не может быть пустым")
	ErrNameEmpty         = errors.New("Имя пользователя не может быть пустым")
	ErrPasswordEmpty     = errors.New("Пароль не может быть пустым")
)

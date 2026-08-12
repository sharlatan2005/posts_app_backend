package serverrors

import "errors"

var (
	ErrWrongPassword = errors.New("Неверно введенный пароль")
)

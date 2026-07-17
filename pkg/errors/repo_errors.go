package errors

import "errors"

var (
	ErrDuplicate = errors.New("duplicate record")
	ErrNotFound  = errors.New("record not found")
	ErrDB        = errors.New("database error")
)

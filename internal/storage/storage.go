package storage

import "errors"

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrAppNotFound   = errors.New("app not found")
	ErrIdeaNotFound  = errors.New("idea not found")
	ErrBoardNotFound = errors.New("board not found")
)

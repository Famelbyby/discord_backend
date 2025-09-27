package storage

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrCantAddBlockedUser = errors.New("can't send friend request to this user")
)

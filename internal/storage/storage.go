package storage

import "errors"

var (
	ErrUserExists             = errors.New("user already exists")
	ErrUserNotFound           = errors.New("user not found")
	ErrCantAddBlockedUser     = errors.New("can't send friend request to this user")
	ErrCantAddAlreadyFriend   = errors.New("can't send friend request to this user, you are already friends")
	ErrUserAlreadyBlocked     = errors.New("this user is already blocked")
	ErrUserAlreadyUnblocked   = errors.New("this user is already unblocked")
	ErrFriendOfferAlreadySent = errors.New("friendship offer  to this user is already sent")
)

package exception

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")

	ErrUserNoPassword  = errors.New("account has no password, probably using social login")
	ErrInvalidPassword = errors.New("invalid password")
)

package user

import "errors"

var (
	ErrNotFound             = errors.New("user not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrWrongPassword        = errors.New("wrong password")
)

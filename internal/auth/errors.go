package auth

import "errors"

var (
	ErrUserExists   = errors.New("user exists")
	ErrInvalidLogin = errors.New("invalid login or password")
)

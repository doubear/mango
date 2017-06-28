package jwt

import (
	"errors"
)

var (
	//ErrTokenLost means token is lost.
	ErrTokenLost = errors.New("token is Lost")

	//ErrTokenInvalid means token is invalid.
	ErrTokenInvalid = errors.New("token is invalid")

	//ErrTokenExpired means token is expired.
	ErrTokenExpired = errors.New("token is expired")
)

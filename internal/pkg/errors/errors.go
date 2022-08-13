package errors

import "errors"

var (
	errInvalidCookie  = errors.New("invalid cookie value")
	errInvalidVersion = errors.New("invalid version")
	errInvalidUint    = errors.New("invalid uint value")
)

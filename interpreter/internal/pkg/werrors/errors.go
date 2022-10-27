package werrors

import "errors"

var (
	ErrInvalidCookie    = errors.New("invalid cookie value")
	ErrInvalidVersion   = errors.New("invalid version")
	ErrInvalidUint      = errors.New("invalid uint value")
	ErrInvalidSectionID = errors.New("invalid section id")
)

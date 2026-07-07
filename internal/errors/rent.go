package errors

import "errors"

var (
	ErrInvalidLat     = errors.New("invalid latitude")
	ErrInvalidLong    = errors.New("invalid longitude")
	ErrInvalidAddress = errors.New("invalid address")
)

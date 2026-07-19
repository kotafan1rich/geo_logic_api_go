package errors

import "errors"

var (
	ErrInvalidSlug   = errors.New("invalid slug")
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidWeight = errors.New("invalid weight")
)

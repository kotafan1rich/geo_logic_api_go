package rent

import "errors"

var (
	ErrRentNotFound      = errors.New("rent not found")
	ErrRentAlreadyExists = errors.New("rent already exists")
)

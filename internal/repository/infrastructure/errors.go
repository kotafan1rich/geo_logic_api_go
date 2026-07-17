package infrastructure

import "errors"

var (
	ErrInfrastructureNotFound          = errors.New("infrastructure not found")
	ErrInfrastructureAlreadyExists     = errors.New("infrastructure already exists")
	ErrInfrastructureTypeNotFound      = errors.New("infrastructure type not found")
	ErrInfrastructureTypeAlreadyExists = errors.New("infrastructure type already exists")
)

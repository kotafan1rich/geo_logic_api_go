package infra

import "errors"

var (
	ErrInfrastructureNotFound          = errors.New("infra not found")
	ErrInfrastructureAlreadyExists     = errors.New("infra already exists")
	ErrInfrastructureTypeNotFound      = errors.New("infra type not found")
	ErrInfrastructureTypeAlreadyExists = errors.New("infra type already exists")
)

package infra

import "errors"

var (
	ErrInfraNotFound          = errors.New("infra not found")
	ErrInfraAlreadyExists     = errors.New("infra already exists")
	ErrInfraTypeNotFound      = errors.New("infra type not found")
	ErrInfraTypeAlreadyExists = errors.New("infra type already exists")
)

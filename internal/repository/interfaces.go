package repository

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetById(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type RentRepository interface {
	Create(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	GetByID(ctx context.Context, id uint64) (*model.Rent, error)
}

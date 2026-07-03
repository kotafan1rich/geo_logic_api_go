package service

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type UserService interface {
	Create(ctx context.Context, tgId uint64) (*model.User, error)
	GetById(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, id uint64, newTgId uint64) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type RentService interface {
	Create(ctx context.Context, lat, long float64, address, info string) (*model.Rent, error)
}

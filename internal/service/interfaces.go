package service

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type UserService interface {
	Create(ctx context.Context, tgId uint64) (*model.User, error)
	GetById(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, id, newTgId uint64) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type RentService interface {
	Create(ctx context.Context, lat, long float64, address string, info *string) (*model.Rent, error)
	GetById(ctx context.Context, id uint64) (*model.Rent, error)
	Update(ctx context.Context, id uint64, lat, long *float64, address, info *string) (*model.Rent, error)
	Delete(ctx context.Context, id uint64) error
	Available(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Rent, error)
}

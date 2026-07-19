package repository

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type RentRepository interface {
	Create(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	GetByID(ctx context.Context, id uint64) (*model.Rent, error)
	Update(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	Delete(ctx context.Context, id uint64) error
	Available(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Rent, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) (*model.Event, error)
	GetByID(ctx context.Context, id uint64) (*model.Event, error)
	Update(ctx context.Context, event *model.Event) (*model.Event, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Event, error)
}

type InfraRepository interface {
	Create(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraObject, error)
	Update(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.InfraObject, error)
}

type InfraTypeRepository interface {
	Create(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraType, error)
	Update(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error)
	Delete(ctx context.Context, id uint64) error
}

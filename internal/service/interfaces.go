package service

import (
	"context"
	"time"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type UserService interface {
	Create(ctx context.Context, tgID uint64) (*model.User, error)
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, id, newTgID uint64) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type RentService interface {
	Create(ctx context.Context, lat, long float64, address string, info *string) (*model.Rent, error)
	GetByID(ctx context.Context, id uint64) (*model.Rent, error)
	Update(ctx context.Context, id uint64, lat, long *float64, address, info *string) (*model.Rent, error)
	Delete(ctx context.Context, id uint64) error
	Available(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Rent, error)
}

type EventService interface {
	Create(ctx context.Context, lat, long float64, date time.Time, info *string) (*model.Event, error)
	GetByID(ctx context.Context, id uint64) (*model.Event, error)
	Update(ctx context.Context, id uint64, lat, long *float64, date *time.Time, info *string) (*model.Event, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Event, error)
}

type InfraTypeService interface {
	Create(ctx context.Context, slug, name string, weight float64, maxRadius uint16) (*model.InfraType, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraType, error)
	Update(ctx context.Context, id uint64, slug, name *string, weight *float64, maxRadius *uint16) (*model.InfraType, error)
	Delete(ctx context.Context, id uint64) error
}

type InfraService interface {
	Create(ctx context.Context, lat, long float64, address string, name *string, infraID uint64) (*model.InfraObject, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraObject, error)
	Update(ctx context.Context, id uint64, lat, long *float64, address, name *string, infraID *uint64) (*model.InfraObject, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error)
}

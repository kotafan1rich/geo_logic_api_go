package infra

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/infra"
)

type InfraRepository interface {
	Create(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraObject, error)
	Update(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error)
}

type InfraService interface {
	Create(ctx context.Context, lat, long float64, address string, name *string, typeID uint64) (*model.InfraObject, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraObject, error)
	Update(ctx context.Context, id uint64, lat, long *float64, address, name *string, typeID *uint64) (*model.InfraObject, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error)
}

type infraService struct {
	repo InfraRepository
	log  logger.Logger
}

func NewInfraService(log logger.Logger, repo InfraRepository) InfraService {
	return &infraService{log: log, repo: repo}
}

func (i *infraService) Create(ctx context.Context, lat, long float64, address string, name *string, infraID uint64) (*model.InfraObject, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		i.log.Warn(
			"invalid lat or long",
			slog.Float64("lat", lat),
			slog.Float64("long", long),
			slog.String("error", err.Error()),
		)
		return nil, apperrors.ValidationError(err.Error())
	}

	newInfra, err := model.NewInfraObject(*geopoint, address, name, model.InfraType{ID: infraID})
	if err != nil {
		return nil, apperrors.ValidationError(err.Error())
	}

	newInfra, err = i.repo.Create(ctx, newInfra)
	if err != nil {
		if errors.Is(err, infra.ErrInfraTypeAlreadyExists) {
			i.log.Warn("infra already exists", slog.String("error", err.Error()))
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		i.log.Error("error creating infra", slog.String("error", err.Error()))
		return nil, err
	}
	return newInfra, nil
}

func (i *infraService) GetByID(ctx context.Context, id uint64) (*model.InfraObject, error) {
	result, err := i.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraNotFound) {
			i.log.Warn("infra not found", slog.String("error", err.Error()))
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		i.log.Error("error getting infra", slog.String("error", err.Error()))
		return nil, err
	}
	return result, nil
}

func (i *infraService) Update(ctx context.Context, id uint64, lat, long *float64, address, name *string, infraID *uint64) (*model.InfraObject, error) {
	oldInfra, err := i.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraNotFound) {
			i.log.Warn(
				"infra not found",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		i.log.Error("error getting infra", slog.String("error", err.Error()))
		return nil, err
	}

	err = oldInfra.Update(lat, long, address, name, infraID)
	if err != nil {
		i.log.Warn("validation error", slog.String("error", err.Error()))
		return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
	}

	updatedInfra, err := i.repo.Update(ctx, oldInfra)
	if err != nil {
		if errors.Is(err, infra.ErrInfraAlreadyExists) {
			i.log.Warn(
				"error updating infra",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		if errors.Is(err, infra.ErrInfraTypeNotFound) {
			i.log.Warn(
				"infra type for update not found",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		i.log.Error("error updating infra", slog.String("error", err.Error()))
		return nil, err
	}
	return updatedInfra, nil
}

func (i *infraService) Delete(ctx context.Context, id uint64) error {
	err := i.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraNotFound) {
			i.log.Warn("infra not found", slog.String("error", err.Error()))
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		i.log.Error(
			"error deleting infra",
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (i *infraService) Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error) {
	results, err := i.repo.Near(ctx, geopoint)
	if err != nil {
		i.log.Error("error getting near events", slog.String("error", err.Error()))
		return nil, err
	}
	return results, nil
}

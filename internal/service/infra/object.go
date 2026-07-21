package infra

import (
	"context"
	"errors"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/infra"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type InfraLogger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type infraService struct {
	repo repository.InfraRepository
	log  InfraLogger
}

func NewInfraService(log InfraLogger, repo repository.InfraRepository) service.InfraService {
	return &infraService{log: log, repo: repo}
}

func (i *infraService) Create(ctx context.Context, lat, long float64, address string, name *string, infraID uint64) (*model.InfraObject, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		i.log.Warn("invalid lat or long", "error", err)
		return nil, apperrors.ValidationError(err.Error())
	}

	newInfra, err := model.NewInfraObject(*geopoint, address, name, model.InfraType{ID: infraID})
	if err != nil {
		return nil, apperrors.ValidationError(err.Error())
	}

	newInfra, err = i.repo.Create(ctx, newInfra)
	if err != nil {
		if errors.Is(err, infra.ErrInfraTypeAlreadyExists) {
			i.log.Warn("infra already exists", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		i.log.Error("error creating infra", "error", err)
		return nil, err
	}
	return newInfra, nil
}

func (i *infraService) GetByID(ctx context.Context, id uint64) (*model.InfraObject, error) {
	result, err := i.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraNotFound) {
			i.log.Warn("infra not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		i.log.Error("error getting infra", "error", err)
		return nil, err
	}
	return result, nil
}

func (i *infraService) Update(ctx context.Context, id uint64, lat, long *float64, address, name *string, infraID *uint64) (*model.InfraObject, error) {
	oldInfra, err := i.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraNotFound) {
			i.log.Warn("infra not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		i.log.Error("error getting infra", "error", err)
		return nil, err
	}

	if lat != nil {
		err = oldInfra.GeoPoint.UpdateLat(*lat)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLat) {
				i.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			i.log.Error("error updating rent", "error", err)
			return nil, err
		}
	}
	if long != nil {
		err = oldInfra.GeoPoint.UpdateLong(*long)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLong) {
				i.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			i.log.Error("error updating infra", "error", err)
			return nil, err
		}
	}

	if address != nil {
		err := oldInfra.UpdateAddress(*address)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLong) {
				i.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			i.log.Error("error updating infra", "error", err)
			return nil, err
		}
	}

	if name != nil {
		oldInfra.UpdateName(name)
	}

	if infraID != nil {
		oldInfra.UpdateTypeID(*infraID)
	}

	updatedInfra, err := i.repo.Update(ctx, oldInfra)
	if err != nil {
		if errors.Is(err, infra.ErrInfraAlreadyExists) {
			i.log.Warn("error updating infra", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		if errors.Is(err, infra.ErrInfraTypeNotFound) {
			i.log.Warn("infra type for update not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		i.log.Error("error updating infra", "error", err)
		return nil, err
	}
	return updatedInfra, nil
}

func (i *infraService) Delete(ctx context.Context, id uint64) error {
	err := i.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraNotFound) {
			i.log.Warn("infra not found", "error", err)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		i.log.Error("error deleting infra", "error", err)
		return err
	}
	return nil
}

func (i *infraService) Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error) {
	results, err := i.repo.Near(ctx, geopoint)
	if err != nil {
		i.log.Error("error getting near events", "error", err)
		return nil, err
	}
	return results, nil
}

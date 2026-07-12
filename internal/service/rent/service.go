package rent

import (
	"context"
	"errors"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent"
)

const (
	defaultRadius = 500
	maxRadius     = 5000
)

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type rentService struct {
	repo repository.RentRepository
	log  Logger
}

func NewRentService(log Logger, repo repository.RentRepository) *rentService {
	return &rentService{repo: repo, log: log}
}

func (s *rentService) Create(ctx context.Context, lat, long float64, address string, info *string) (*model.Rent, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidLat) || errors.Is(err, apperrors.ErrInvalidLong) {
			s.log.Warn("validation error", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
		s.log.Error("error creating rent", "error", err)
		return nil, err
	}
	rentToCreate := model.NewRent(geopoint, address, info)

	createdRent, err := s.repo.Create(ctx, rentToCreate)
	if err != nil {
		if errors.Is(err, rent.ErrRentAlreadyExists) {
			s.log.Warn("error creating rent", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error("error creating rent", "error", err)
		return nil, err
	}
	return createdRent, nil
}

func (s *rentService) GetById(ctx context.Context, id uint64) (*model.Rent, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			s.log.Warn("rent not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error getting rent", "error", err)
		return nil, err
	}
	return result, nil
}

func (s *rentService) Update(ctx context.Context, id uint64, lat, long *float64, address, info *string) (*model.Rent, error) {
	oldRent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			s.log.Warn("rent not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error getting rent", "error", err)
		return nil, err
	}

	if lat != nil {
		err = oldRent.Geopoint.UpdateLat(*lat)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLat) {
				s.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			s.log.Error("error updating rent", "error", err)
			return nil, err
		}
	}
	if long != nil {
		err = oldRent.Geopoint.UpdateLong(*long)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLong) {
				s.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			s.log.Error("error updating rent", "error", err)
			return nil, err
		}
	}

	if address != nil {
		err = oldRent.UpdateAddress(*address)
		if err != nil {
			s.log.Warn("validation error", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
	}

	if info != nil {
		oldRent.UpdateInfo(*info)
	}

	updatedRent, err := s.repo.Update(ctx, oldRent)
	if err != nil {
		if errors.Is(err, rent.ErrRentAlreadyExists) {
			s.log.Warn("error updating rent", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error("error updating rent", "error", err)
		return nil, err
	}
	return updatedRent, nil
}

func (s *rentService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			s.log.Warn("rent not found", "error", err)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error deleting rent", "error", err)
		return err
	}
	return nil
}

func (s *rentService) Available(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Rent, error) {
	var finalRadius uint16
	if radius == nil {
		finalRadius = defaultRadius
	} else {
		finalRadius = *radius
		if finalRadius == 0 || finalRadius > maxRadius {
			s.log.Warn("invalid radius", "error", apperrors.ErrInvalidRadius)
			return nil, apperrors.Wrap(apperrors.ErrInvalidRadius, apperrors.ValidationError(apperrors.ErrInvalidRadius.Error()))
		}
	}
	results, err := s.repo.Available(ctx, geopoint, finalRadius)
	if err != nil {
		s.log.Error("error getting available rents", "error", err)
		return nil, err
	}
	return results, nil
}

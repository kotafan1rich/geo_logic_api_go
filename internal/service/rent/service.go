package rent

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent"
)

const (
	defaultRadius = 500
	maxRadius     = 5000
)

type RentRepository interface {
	Create(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	GetByID(ctx context.Context, id uint64) (*model.Rent, error)
	Update(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	Delete(ctx context.Context, id uint64) error
	Available(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Rent, error)
}

type RentService struct {
	repo RentRepository
	log  logger.Logger
}

func NewRentService(log logger.Logger, repo RentRepository) *RentService {
	return &RentService{repo: repo, log: log}
}

func (s *RentService) Create(ctx context.Context, lat, long float64, address string, info *string) (*model.Rent, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidLat) || errors.Is(err, apperrors.ErrInvalidLong) {
			s.log.Warn(
				"validation error",
				slog.Float64("lat", lat),
				slog.Float64("long", long),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
		s.log.Error(
			"error creating rent",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	rentToCreate := model.NewRent(geopoint, address, info)

	createdRent, err := s.repo.Create(ctx, rentToCreate)
	if err != nil {
		if errors.Is(err, rent.ErrRentAlreadyExists) {
			s.log.Warn(
				"error creating rent",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error(
			"error creating rent",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return createdRent, nil
}

func (s *RentService) GetByID(ctx context.Context, id uint64) (*model.Rent, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			s.log.Warn(
				"rent not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting rent",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return result, nil
}

func (s *RentService) Update(ctx context.Context, id uint64, lat, long *float64, address, info *string) (*model.Rent, error) {
	oldRent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			s.log.Warn(
				"rent not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting rent",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	err = oldRent.Update(lat, long, address, info)
	if err != nil {
		s.log.Warn(
			"validation error",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
	}

	updatedRent, err := s.repo.Update(ctx, oldRent)
	if err != nil {
		if errors.Is(err, rent.ErrRentAlreadyExists) {
			s.log.Warn(
				"error updating rent",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error(
			"error updating rent",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return updatedRent, nil
}

func (s *RentService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			s.log.Warn(
				"rent not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error deleting rent",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (s *RentService) Available(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Rent, error) {
	var finalRadius uint16
	if radius == nil {
		finalRadius = defaultRadius
	} else {
		finalRadius = *radius
		if finalRadius == 0 || finalRadius > maxRadius {
			s.log.Warn(
				"invalid radius",
				slog.Uint64("radius", uint64(finalRadius)),
				slog.String("error", apperrors.ErrInvalidRadius.Error()),
			)
			return nil, apperrors.Wrap(apperrors.ErrInvalidRadius, apperrors.ValidationError(apperrors.ErrInvalidRadius.Error()))
		}
	}
	results, err := s.repo.Available(ctx, geopoint, finalRadius)
	if err != nil {
		s.log.Error(
			"error getting available rents",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return results, nil
}

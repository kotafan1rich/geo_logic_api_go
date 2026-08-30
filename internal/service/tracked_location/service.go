package trackedlocation

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/tracked_location"
)

type TrackedLocationRepository interface {
	Create(ctx context.Context, location *model.TrackedLocation) (*model.TrackedLocation, error)
	GetByID(ctx context.Context, id uint64) (*model.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uint64) ([]model.TrackedLocation, error)
	Update(ctx context.Context, location *model.TrackedLocation) (*model.TrackedLocation, error)
	Delete(ctx context.Context, id uint64) error
}

type TrackedLocationService interface {
	Create(ctx context.Context, userID uint64, lat, long float64) (*model.TrackedLocation, error)
	GetByID(ctx context.Context, id uint64) (*model.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uint64) ([]model.TrackedLocation, error)
	Update(ctx context.Context, id uint64, userID *uint64, lat, long *float64) (*model.TrackedLocation, error)
	Delete(ctx context.Context, id uint64) error
}

type locationService struct {
	repo TrackedLocationRepository
	log  logger.Logger
}

func NewTrackedLocationService(log logger.Logger, repo TrackedLocationRepository) TrackedLocationService {
	return &locationService{repo: repo, log: log}
}

func (s *locationService) Create(ctx context.Context, userID uint64, lat, long float64) (*model.TrackedLocation, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidLat) || errors.Is(err, apperrors.ErrInvalidLong) {
			s.log.Warn(
				"validation error",
				slog.Uint64("userID", userID),
				slog.Float64("lat", lat),
				slog.Float64("long", long),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
		s.log.Error(
			"error creating tracked location",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	locationToCreate := model.NewTrackedLocation(userID, geopoint)

	createdLocation, err := s.repo.Create(ctx, locationToCreate)
	if err != nil {
		if errors.Is(err, trackedlocation.ErrTrackedLocationAlreadyExists) {
			s.log.Warn(
				"error creating tracked location",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error(
			"error creating tracked location",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return createdLocation, nil
}

func (s *locationService) GetByID(ctx context.Context, id uint64) (*model.TrackedLocation, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, trackedlocation.ErrTrackedLocationNotFound) {
			s.log.Warn(
				"tracked location not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting tracked location",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return result, nil
}

func (s *locationService) GetByUserID(ctx context.Context, userID uint64) ([]model.TrackedLocation, error) {
	result, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, trackedlocation.ErrTrackedLocationNotFound) {
			s.log.Warn(
				"tracked location not found",
				slog.Uint64("userID", userID),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting tracked location",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return result, nil
}

func (s *locationService) Update(ctx context.Context, id uint64, userID *uint64, lat, long *float64) (*model.TrackedLocation, error) {
	oldLocation, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, trackedlocation.ErrTrackedLocationNotFound) {
			s.log.Warn(
				"tracked location not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting tracked location",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	err = oldLocation.Update(userID, lat, long)
	if err != nil {
		s.log.Warn(
			"validation error",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
	}

	updatedLocation, err := s.repo.Update(ctx, oldLocation)
	if err != nil {
		if errors.Is(err, trackedlocation.ErrTrackedLocationAlreadyExists) {
			s.log.Warn(
				"error updating tracked location",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error(
			"error updating tracked location",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return updatedLocation, nil
}

func (s *locationService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, trackedlocation.ErrTrackedLocationNotFound) {
			s.log.Warn(
				"tracked location not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error deleting tracked location",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

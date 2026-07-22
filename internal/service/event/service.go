package event

import (
	"context"
	"errors"
	"log/slog"
	"time"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/event"
)

const (
	defaultRadius = 500
	maxRadius     = 5000
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) (*model.Event, error)
	GetByID(ctx context.Context, id uint64) (*model.Event, error)
	Update(ctx context.Context, event *model.Event) (*model.Event, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Event, error)
}

type EventService interface {
	Create(ctx context.Context, lat, long float64, date time.Time, info *string) (*model.Event, error)
	GetByID(ctx context.Context, id uint64) (*model.Event, error)
	Update(ctx context.Context, id uint64, lat, long *float64, date *time.Time, info *string) (*model.Event, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Event, error)
}

type eventService struct {
	repo EventRepository
	log  logger.Logger
}

func NewEventService(log logger.Logger, repo EventRepository) EventService {
	return &eventService{repo: repo, log: log}
}

func (s *eventService) Create(ctx context.Context, lat, long float64, date time.Time, info *string) (*model.Event, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidLat) || errors.Is(err, apperrors.ErrInvalidLong) {
			s.log.Warn(
				"validation error",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
		s.log.Error(
			"error creating event",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	eventToCreate := model.NewEvent(*geopoint, date, info)

	createdEvent, err := s.repo.Create(ctx, eventToCreate)
	if err != nil {
		if errors.Is(err, event.ErrEventAlreadyExists) {
			s.log.Warn(
				"event already exists",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error(
			"error creating event",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return createdEvent, nil
}

func (s *eventService) GetByID(ctx context.Context, id uint64) (*model.Event, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			s.log.Warn(
				"event not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting event",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return result, nil
}

func (s *eventService) Update(ctx context.Context, id uint64, lat, long *float64, date *time.Time, info *string) (*model.Event, error) {
	oldEvent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			s.log.Warn(
				"event not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error getting event",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	err = oldEvent.Update(lat, long, date, info)
	if err != nil {
		s.log.Warn(
			"validation error",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
	}

	updatedEvent, err := s.repo.Update(ctx, oldEvent)
	if err != nil {
		if errors.Is(err, event.ErrEventAlreadyExists) {
			s.log.Warn(
				"error updating event",
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error(
			"error updating event",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return updatedEvent, nil
}

func (s *eventService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			s.log.Warn(
				"event not found",
				slog.String("error", err.Error()),
			)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error(
			"error deleting event",
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

func (s *eventService) Near(ctx context.Context, geopoint *model.GeoPoint, radius *uint16) ([]model.Event, error) {
	var finalRadius uint16
	if radius == nil {
		finalRadius = defaultRadius
	} else {
		finalRadius = *radius
		if finalRadius == 0 || finalRadius > maxRadius {
			s.log.Warn(
				"invalid radius",
				slog.String("error", apperrors.ErrInvalidRadius.Error()),
			)
			return nil, apperrors.Wrap(apperrors.ErrInvalidRadius, apperrors.ValidationError(apperrors.ErrInvalidRadius.Error()))
		}
	}
	results, err := s.repo.Near(ctx, geopoint, finalRadius)
	if err != nil {
		s.log.Error(
			"error getting near events",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return results, nil
}

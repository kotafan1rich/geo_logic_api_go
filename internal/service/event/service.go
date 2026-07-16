package event

import (
	"context"
	"errors"
	"time"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/event"
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

type eventService struct {
	repo repository.EventRepository
	log  Logger
}

func NewEventService(log Logger, repo repository.EventRepository) *eventService {
	return &eventService{repo: repo, log: log}
}

func (s *eventService) Create(ctx context.Context, lat, long float64, date time.Time, info *string) (*model.Event, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidLat) || errors.Is(err, apperrors.ErrInvalidLong) {
			s.log.Warn("validation error", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
		s.log.Error("error creating event", "error", err)
		return nil, err
	}
	eventToCreate := model.NewEvent(*geopoint, date, info)

	createdEvent, err := s.repo.Create(ctx, eventToCreate)
	if err != nil {
		if errors.Is(err, event.ErrEventAlreadyExists) {
			s.log.Warn("error creating event", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error("error creating event", "error", err)
		return nil, err
	}
	return createdEvent, nil
}

func (s *eventService) GetById(ctx context.Context, id uint64) (*model.Event, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			s.log.Warn("event not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error getting event", "error", err)
		return nil, err
	}
	return result, nil
}

func (s *eventService) Update(ctx context.Context, id uint64, lat, long *float64, date *time.Time, info *string) (*model.Event, error) {
	oldEvent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			s.log.Warn("event not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error getting event", "error", err)
		return nil, err
	}

	if lat != nil {
		err = oldEvent.Geopint.UpdateLat(*lat)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLat) {
				s.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			s.log.Error("error updating event", "error", err)
			return nil, err
		}
	}
	if long != nil {
		err = oldEvent.Geopint.UpdateLong(*long)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLong) {
				s.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			s.log.Error("error updating event", "error", err)
			return nil, err
		}
	}

	if date != nil {
		oldEvent.UpdateDate(*date)
	}

	if info != nil {
		oldEvent.UpdateInfo(info)
	}

	updatedEvent, err := s.repo.Update(ctx, oldEvent)
	if err != nil {
		if errors.Is(err, event.ErrEventAlreadyExists) {
			s.log.Warn("error updating event", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error("error updating event", "error", err)
		return nil, err
	}
	return updatedEvent, nil
}

func (s *eventService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, event.ErrEventNotFound) {
			s.log.Warn("event not found", "error", err)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error deleting event", "error", err)
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
			s.log.Warn("invalid radius", "error", apperrors.ErrInvalidRadius)
			return nil, apperrors.Wrap(apperrors.ErrInvalidRadius, apperrors.ValidationError(apperrors.ErrInvalidRadius.Error()))
		}
	}
	results, err := s.repo.Near(ctx, geopoint, finalRadius)
	if err != nil {
		s.log.Error("error getting near events", "error", err)
		return nil, err
	}
	return results, nil
}

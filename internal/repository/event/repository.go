package event

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/event/dbmodel"
	"gorm.io/gorm"
)

type eventRepository struct {
	db database.DB
}

func NewRepository(db database.DB) *eventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) (*model.Event, error) {
	eventModel := dbmodel.ToEventModel(event)
	err := r.db.GORM().WithContext(ctx).Create(eventModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrEventAlreadyExists
		}
		return nil, err
	}
	return dbmodel.ToEvent(eventModel), nil
}

func (r *eventRepository) GetByID(ctx context.Context, id uint64) (*model.Event, error) {
	var event dbmodel.Event
	err := r.db.GORM().WithContext(ctx).First(&event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}
	return dbmodel.ToEvent(&event), err
}

func (r *eventRepository) Update(ctx context.Context, event *model.Event) (*model.Event, error) {
	eventModel := dbmodel.ToEventModel(event)
	result := r.db.GORM().WithContext(ctx).Model(&dbmodel.Event{}).Where("id = ?", eventModel.ID).Updates(map[string]any{
		"date":     eventModel.Date,
		"info":     eventModel.Info,
		"location": gorm.Expr("ST_SetSRID(ST_Point(?, ?), 4326)", eventModel.Location.Long, eventModel.Location.Lat),
	})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrEventAlreadyExists
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrEventNotFound
	}
	return event, nil
}

func (r *eventRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.GORM().WithContext(ctx).Delete(&dbmodel.Event{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrEventNotFound
	}
	return nil
}

func (r *eventRepository) Near(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Event, error) {
	var events []dbmodel.Event

	err := r.db.GORM().WithContext(ctx).Model(&dbmodel.Event{}).Where(
		"ST_DWithin(location, ST_SetSRID(ST_Point(?, ?), 4326)::geography, ?)",
		geopoint.Long,
		geopoint.Lat,
		radius,
	).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return dbmodel.ToEventSlice(events), nil
}

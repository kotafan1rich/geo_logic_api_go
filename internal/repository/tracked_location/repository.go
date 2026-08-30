package trackedlocation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/tracked_location/dbmodel"
	"gorm.io/gorm"
)

type TrackedLocationRepository interface {
	Create(ctx context.Context, location *model.TrackedLocation) (*model.TrackedLocation, error)
	GetByID(ctx context.Context, id uint64) (*model.TrackedLocation, error)
	GetByUserID(ctx context.Context, userID uint64) ([]model.TrackedLocation, error)
	Update(ctx context.Context, location *model.TrackedLocation) (*model.TrackedLocation, error)
	Delete(ctx context.Context, id uint64) error
}

type locationRepository struct {
	db database.DB
}

func NewRepository(db database.DB) TrackedLocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(ctx context.Context, location *model.TrackedLocation) (*model.TrackedLocation, error) {
	locationModel := dbmodel.ToModel(location)
	err := r.db.GORM().WithContext(ctx).Create(locationModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrTrackedLocationAlreadyExists
		}
		return nil, err
	}
	return dbmodel.ToDomain(locationModel), nil
}

func (r *locationRepository) GetByID(ctx context.Context, id uint64) (*model.TrackedLocation, error) {
	var location dbmodel.TrackedLocation
	err := r.db.GORM().WithContext(ctx).First(&location, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTrackedLocationNotFound
		}
		return nil, err
	}
	return dbmodel.ToDomain(&location), nil
}

func (r *locationRepository) GetByUserID(ctx context.Context, userID uint64) ([]model.TrackedLocation, error) {
	var locations []dbmodel.TrackedLocation
	err := r.db.GORM().WithContext(ctx).Model(&dbmodel.TrackedLocation{}).Where(
		"user_id = ?",
		userID,
	).Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return dbmodel.ToDomainSlice(locations), nil
}

func (r *locationRepository) Update(ctx context.Context, location *model.TrackedLocation) (*model.TrackedLocation, error) {
	locationModel := dbmodel.ToModel(location)
	result := r.db.GORM().WithContext(ctx).Model(&dbmodel.TrackedLocation{}).Where("id = ?", locationModel.ID).Updates(map[string]any{
		"user_id":  locationModel.UserID,
		"location": gorm.Expr("ST_SetSRID(ST_Point(?, ?), 4326)", locationModel.Location.Long, locationModel.Location.Lat),
	})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrTrackedLocationAlreadyExists
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrTrackedLocationNotFound
	}
	return location, nil
}

func (r *locationRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.GORM().WithContext(ctx).Delete(&dbmodel.TrackedLocation{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTrackedLocationNotFound
	}
	return nil
}

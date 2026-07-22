package rent

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent/dbmodel"
	"gorm.io/gorm"
)

type RentRepository interface {
	Create(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	GetByID(ctx context.Context, id uint64) (*model.Rent, error)
	Update(ctx context.Context, rent *model.Rent) (*model.Rent, error)
	Delete(ctx context.Context, id uint64) error
	Available(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Rent, error)
}

type rentRepository struct {
	db database.DB
}

func NewRepository(db database.DB) RentRepository {
	return &rentRepository{db: db}
}

func (r *rentRepository) Create(ctx context.Context, rent *model.Rent) (*model.Rent, error) {
	rentModel := dbmodel.ToRentModel(rent)
	err := r.db.GORM().WithContext(ctx).Create(rentModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrRentAlreadyExists
		}
		return nil, err
	}
	return dbmodel.ToRent(rentModel), nil
}

func (r *rentRepository) GetByID(ctx context.Context, id uint64) (*model.Rent, error) {
	var rent dbmodel.Rent
	err := r.db.GORM().WithContext(ctx).First(&rent, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRentNotFound
		}
		return nil, err
	}
	return dbmodel.ToRent(&rent), err
}

func (r *rentRepository) Update(ctx context.Context, rent *model.Rent) (*model.Rent, error) {
	rentModel := dbmodel.ToRentModel(rent)
	result := r.db.GORM().WithContext(ctx).Model(&dbmodel.Rent{}).Where("id = ?", rentModel.ID).Updates(map[string]any{
		"address":  rentModel.Address,
		"info":     rentModel.Info,
		"location": gorm.Expr("ST_SetSRID(ST_Point(?, ?), 4326)", rentModel.Location.Long, rentModel.Location.Lat),
	})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrRentAlreadyExists
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrRentNotFound
	}
	return rent, nil
}

func (r *rentRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.GORM().WithContext(ctx).Delete(&dbmodel.Rent{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrRentNotFound
	}
	return nil
}

func (r *rentRepository) Available(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.Rent, error) {
	var rents []dbmodel.Rent

	err := r.db.GORM().WithContext(ctx).Model(&dbmodel.Rent{}).Where(
		"ST_DWithin(location, ST_SetSRID(ST_Point(?, ?), 4326)::geography, ?)",
		geopoint.Long,
		geopoint.Lat,
		radius,
	).Find(&rents).Error
	if err != nil {
		return nil, err
	}
	return dbmodel.ToRentSlice(rents), nil
}

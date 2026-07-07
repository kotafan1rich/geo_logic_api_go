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

type rentRepository struct {
	db database.DB
}

func NewRepository(db database.DB) *rentRepository {
	return &rentRepository{db: db}
}

func (r *rentRepository) Create(ctx context.Context, rent *model.Rent) (*model.Rent, error) {
	point := dbmodel.DBGeoPoint(*rent.Geopoint)
	wktPoint, err := point.Value()
	if err != nil {
		return nil, err
	}

	var id uint64
	err = r.db.GORM().WithContext(ctx).Raw(
		`INSERT INTO "rents" ("location", "address", "info") 
		 VALUES (ST_GeomFromText(?, 4326), ?, ?)
		 RETURNING id`,
		wktPoint, rent.Address, rent.Info,
	).Scan(&id).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrRentAlreadyExists
		}
		return nil, err
	}
	rent.ID = id
	return rent, nil
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
	result := r.db.GORM().WithContext(ctx).Model(&rentModel).Updates(&rentModel)
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

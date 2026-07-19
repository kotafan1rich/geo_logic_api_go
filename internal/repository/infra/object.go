package infra

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/infra/dbmodel"
	"gorm.io/gorm"
)

type infraRepository struct {
	db database.DB
}

func NewInfraRepository(db database.DB) *infraRepository {
	return &infraRepository{db: db}
}

func (r *infraRepository) Create(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error) {
	point := database.DBGeoPoint(infra.GeoPoint)
	wktPoint, err := point.Value()
	if err != nil {
		return nil, err
	}

	var id uint64
	err = r.db.GORM().WithContext(ctx).Raw(
		`INSERT INTO "structure_objects" (
			"location",
			"address",
			"name",
			"type_id"
		)
		 VALUES (ST_GeomFromText(?, 4326), ?, ?, ?)
		 RETURNING id`,
		wktPoint, infra.Address, infra.Name, infra.Type.ID,
	).Scan(&id).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case database.ErrPgUniqueViolation:
				return nil, ErrInfraAlreadyExists
			case database.PgErrForeignKeyViolation:
				return nil, ErrInfraTypeNotFound
			}
		}
		return nil, err
	}
	infra.ID = id
	return infra, nil
}

func (r *infraRepository) GetByID(ctx context.Context, id uint64) (*model.InfraObject, error) {
	var infra dbmodel.InfraObject
	err := r.db.GORM().WithContext(ctx).First(&infra, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInfraNotFound
		}
		return nil, err
	}
	return dbmodel.ToObject(&infra), err
}

func (r *infraRepository) Update(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error) {
	infraModel := dbmodel.ToObjectModel(infra)
	result := r.db.GORM().WithContext(ctx).Model(&dbmodel.InfraObject{}).Where("id = ?", infraModel.ID).Updates(map[string]any{
		"address":  infraModel.Address,
		"name":     infraModel.Name,
		"type_id":  infraModel.TypeID,
		"location": gorm.Expr("ST_SetSRID(ST_Point(?, ?), 4326)", infraModel.Location.Long, infraModel.Location.Lat),
	})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			switch pgErr.Code {
			case database.ErrPgUniqueViolation:
				return nil, ErrInfraAlreadyExists
			case database.PgErrForeignKeyViolation:
				return nil, ErrInfraTypeNotFound
			}
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrInfraNotFound
	}
	return infra, nil
}

func (r *infraRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.GORM().WithContext(ctx).Delete(&dbmodel.InfraObject{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInfraNotFound
	}
	return nil
}

func (r *infraRepository) Near(ctx context.Context, geopoint *model.GeoPoint, radius uint16) ([]model.InfraObject, error) {
	var infras []dbmodel.InfraObject

	err := r.db.GORM().WithContext(ctx).Model(&dbmodel.InfraObject{}).Where(
		"ST_DWithin(location, ST_SetSRID(ST_Point(?, ?), 4326)::geography, ?)",
		geopoint.Long,
		geopoint.Lat,
		radius,
	).Find(&infras).Error
	if err != nil {
		return nil, err
	}
	return dbmodel.ToObjectSlice(infras), nil
}

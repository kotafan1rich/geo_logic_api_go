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

type InfraRepository interface {
	Create(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraObject, error)
	Update(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error)
	Delete(ctx context.Context, id uint64) error
	Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error)
}

type infraRepository struct {
	db database.DB
}

func NewInfraRepository(db database.DB) InfraRepository {
	return &infraRepository{db: db}
}

func (r *infraRepository) Create(ctx context.Context, infra *model.InfraObject) (*model.InfraObject, error) {
	dbInfra := dbmodel.ToObjectModel(infra)
	err := r.db.GORM().WithContext(ctx).Create(&dbInfra).Error
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case database.ErrPgUniqueViolation:
				return nil, ErrInfraAlreadyExists
			case database.PgErrForeignKeyViolation:
				return nil, ErrInfraTypeNotFound
			}
		}
		return nil, err
	}

	err = r.db.GORM().WithContext(ctx).Preload("Type").First(&dbInfra, dbInfra.ID).Error
	if err != nil {
		return nil, err
	}
	return dbmodel.ToObject(dbInfra), nil
}

func (r *infraRepository) GetByID(ctx context.Context, id uint64) (*model.InfraObject, error) {
	var infra dbmodel.InfraObject
	err := r.db.GORM().WithContext(ctx).Preload("Type").First(&infra, id).Error
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
		if pgErr, ok := errors.AsType[*pgconn.PgError](result.Error); ok {
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

func (r *infraRepository) Near(ctx context.Context, geopoint *model.GeoPoint) ([]model.InfraObject, error) {
	var infras []dbmodel.InfraObject

	err := r.db.GORM().WithContext(ctx).
		Preload("Type").
		Joins("JOIN infra_types t ON infra_objects.type_id = t.id").
		Where(`ST_DWithin(
			ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography,
			infra_objects.location::geography,
			t.max_radius
		)`, geopoint.Long, geopoint.Lat).
		Find(&infras).Error
	if err != nil {
		return nil, err
	}
	return dbmodel.ToObjectSlice(infras), nil
}

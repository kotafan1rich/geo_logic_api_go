package infrastructure

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/infrastructure/dbmodel"
	"gorm.io/gorm"
)

type infrastructureTypeRepository struct {
	db database.DB
}

func NewInfrastructureTypeRepository(db database.DB) *infrastructureTypeRepository {
	return &infrastructureTypeRepository{db: db}
}

func (t *infrastructureTypeRepository) Create(ctx context.Context, infraType *model.InfrastructureType) (*model.InfrastructureType, error) {
	typeModel := dbmodel.ToTypeModel(infraType)
	err := t.db.GORM().WithContext(ctx).Create(&typeModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrInfrastructureTypeAlreadyExists
		}
		return nil, err
	}
	infraType.ID = typeModel.ID
	return infraType, nil
}

func (t *infrastructureTypeRepository) GetById(ctx context.Context, id uint64) (*model.InfrastructureType, error) {
	var typeModel dbmodel.InfrastructureType
	err := t.db.GORM().WithContext(ctx).First(&typeModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInfrastructureTypeNotFound
		}
		return nil, err
	}

	return dbmodel.ToType(&typeModel), nil
}

func (t *infrastructureTypeRepository) Update(ctx context.Context, infraType *model.InfrastructureType) (*model.InfrastructureType, error) {
	typeModel := dbmodel.ToTypeModel(infraType)
	result := t.db.GORM().WithContext(ctx).Model(&dbmodel.InfrastructureType{}).Where("id = ?", infraType.ID).Updates(map[string]any{
		"slug":       typeModel.Slug,
		"name":       typeModel.Name,
		"weight":     typeModel.Weight,
		"max_radius": typeModel.MaxRadius,
	})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrInfrastructureTypeAlreadyExists
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrInfrastructureTypeNotFound
	}
	return infraType, nil
}

func (t *infrastructureTypeRepository) Delete(ctx context.Context, id uint64) error {
	result := t.db.GORM().WithContext(ctx).Delete(&dbmodel.InfrastructureType{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInfrastructureTypeNotFound
	}
	return nil
}

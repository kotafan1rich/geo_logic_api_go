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

type InfraTypeRepository interface {
	Create(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraType, error)
	Update(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error)
	Delete(ctx context.Context, id uint64) error
}

type infraTypeRepository struct {
	db database.DB
}

func NewInfraTypeRepository(db database.DB) InfraTypeRepository {
	return &infraTypeRepository{db: db}
}

func (t *infraTypeRepository) Create(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error) {
	typeModel := dbmodel.ToTypeModel(infraType)
	err := t.db.GORM().WithContext(ctx).Create(&typeModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrInfraTypeAlreadyExists
		}
		return nil, err
	}
	infraType.ID = typeModel.ID
	return infraType, nil
}

func (t *infraTypeRepository) GetByID(ctx context.Context, id uint64) (*model.InfraType, error) {
	var typeModel dbmodel.InfraType
	err := t.db.GORM().WithContext(ctx).First(&typeModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInfraTypeNotFound
		}
		return nil, err
	}

	return dbmodel.ToType(&typeModel), nil
}

func (t *infraTypeRepository) Update(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error) {
	typeModel := dbmodel.ToTypeModel(infraType)
	result := t.db.GORM().WithContext(ctx).Model(&dbmodel.InfraType{}).Where("id = ?", infraType.ID).Updates(map[string]any{
		"slug":       typeModel.Slug,
		"name":       typeModel.Name,
		"weight":     typeModel.Weight,
		"max_radius": typeModel.MaxRadius,
	})
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrInfraTypeAlreadyExists
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrInfraTypeNotFound
	}
	return infraType, nil
}

func (t *infraTypeRepository) Delete(ctx context.Context, id uint64) error {
	result := t.db.GORM().WithContext(ctx).Delete(&dbmodel.InfraType{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInfraTypeNotFound
	}
	return nil
}

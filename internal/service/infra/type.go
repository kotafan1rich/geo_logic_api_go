package infra

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/infra"
)

type InfraTypeRepository interface {
	Create(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error)
	GetByID(ctx context.Context, id uint64) (*model.InfraType, error)
	Update(ctx context.Context, infraType *model.InfraType) (*model.InfraType, error)
	Delete(ctx context.Context, id uint64) error
}

type TypeService struct {
	repo InfraTypeRepository
	log  logger.Logger
}

func NewTypeService(log logger.Logger, repo InfraTypeRepository) *TypeService {
	return &TypeService{log: log, repo: repo}
}

func (t *TypeService) Create(ctx context.Context, slug, name string, weight float64, maxRadius uint16) (*model.InfraType, error) {
	newType, err := model.NewInfraType(slug, name, weight, maxRadius)
	if err != nil {
		return nil, apperrors.ValidationError(err.Error())
	}

	newType, err = t.repo.Create(ctx, newType)
	if err != nil {
		if errors.Is(err, infra.ErrInfraTypeAlreadyExists) {
			t.log.Warn(
				"type already exists",
				slog.String("slug", slug),
				slog.String("name", name),
				slog.Float64("weight", weight),
				slog.Uint64("maxRadius", uint64(maxRadius)),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		t.log.Error(
			"error creating type",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return newType, nil
}

func (t *TypeService) GetByID(ctx context.Context, id uint64) (*model.InfraType, error) {
	result, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraTypeNotFound) {
			t.log.Warn(
				"type not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		t.log.Error(
			"error getting type",
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return result, nil
}

func (t *TypeService) Update(ctx context.Context, id uint64, slug, name *string, weight *float64, maxRadius *uint16) (*model.InfraType, error) {
	oldType, err := t.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraTypeNotFound) {
			t.log.Warn(
				"type not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		t.log.Error(
			"error getting type",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	err = oldType.Update(slug, name, weight, maxRadius)
	if err != nil {
		t.log.Warn(
			"validation error",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
	}

	updatedType, err := t.repo.Update(ctx, oldType)
	if err != nil {
		if errors.Is(err, infra.ErrInfraAlreadyExists) {
			t.log.Warn(
				"error updating type",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		t.log.Error(
			"error updating type",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return updatedType, nil
}

func (t *TypeService) Delete(ctx context.Context, id uint64) error {
	err := t.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfraTypeNotFound) {
			t.log.Warn(
				"type not found",
				slog.Uint64("id", id),
				slog.String("error", err.Error()),
			)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		t.log.Error(
			"error deleting type",
			slog.Uint64("id", id),
			slog.String("error", err.Error()),
		)
		return err
	}
	return nil
}

package infra

import (
	"context"
	"errors"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/infra"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type typeService struct {
	repo repository.InfrastructureTypeRepository
	log  Logger
}

func NewTypeService(log Logger, repo repository.InfrastructureTypeRepository) service.InfrastructureTypeService {
	return &typeService{log: log, repo: repo}
}

func (t *typeService) Create(ctx context.Context, slug, name string, weight float64, maxRadius uint16) (*model.InfrastructureType, error) {
	newType, err := model.NewInfrastructureType(slug, name, weight, maxRadius)
	if err != nil {
		return nil, apperrors.ValidationError(err.Error())
	}

	newType, err = t.repo.Create(ctx, newType)
	if err != nil {
		if errors.Is(err, infra.ErrInfrastructureTypeAlreadyExists) {
			t.log.Warn("type already exists", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		t.log.Error("error creating type", "error", err)
		return nil, err
	}
	return newType, nil
}

func (t *typeService) GetById(ctx context.Context, id uint64) (*model.InfrastructureType, error) {
	result, err := t.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfrastructureTypeNotFound) {
			t.log.Warn("type not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		t.log.Error("error getting type", "error", err)
		return nil, err
	}
	return result, nil
}

func (t *typeService) Update(ctx context.Context, id uint64, slug, name *string, weight *float64, maxRadius *uint16) (*model.InfrastructureType, error) {
	oldType, err := t.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfrastructureTypeNotFound) {
			t.log.Warn("type not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		t.log.Error("error getting type", "error", err)
		return nil, err
	}

	if slug != nil {
		err = oldType.UpdateSlug(*slug)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidSlug) {
				t.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			t.log.Error("error updating type", "error", err)
			return nil, err
		}
	}

	if name != nil {
		err = oldType.UpdateName(*name)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidName) {
				t.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			t.log.Error("error updating type", "error", err)
			return nil, err
		}
	}

	if weight != nil {
		err = oldType.UpdateWeight(*weight)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidWeight) {
				t.log.Warn("validation error", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			t.log.Error("error updating type", "error", err)
			return nil, err
		}
	}

	if maxRadius != nil {
		oldType.UpdateMaxRadius(*maxRadius)
	}

	updatedType, err := t.repo.Update(ctx, oldType)
	if err != nil {
		if errors.Is(err, infra.ErrInfrastructureAlreadyExists) {
			t.log.Warn("error updating type", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		t.log.Error("error updating type", "error", err)
		return nil, err
	}
	return updatedType, nil
}

func (t *typeService) Delete(ctx context.Context, id uint64) error {
	err := t.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, infra.ErrInfrastructureTypeNotFound) {
			t.log.Warn("type not found", "error", err)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		t.log.Error("error deleting type", "error", err)
		return err
	}
	return nil
}

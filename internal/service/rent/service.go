package rent

import (
	"context"
	"errors"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/rent"
)

type Logger interface {
	Error(msg string, args ...any)
}

type rentService struct {
	repo repository.RentRepository
	log  Logger
}

func NewRentService(log Logger, repo repository.RentRepository) *rentService {
	return &rentService{repo: repo, log: log}
}

func (s *rentService) Create(ctx context.Context, lat, long float64, address, info string) (*model.Rent, error) {
	geopoint, err := model.NewGeoPoint(lat, long)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidLat) || errors.Is(err, apperrors.ErrInvalidLong) {
			return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
		}
		return nil, err
	}
	rentToCreate := model.NewRent(geopoint, address, info)

	createdRent, err := s.repo.Create(ctx, rentToCreate)
	if err != nil {
		if errors.Is(err, rent.ErrRentAlreadyExists) {
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		return nil, err
	}

	return createdRent, nil
}

func (s *rentService) GetById(ctx context.Context, id uint64) (*model.Rent, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		return nil, err
	}
	return result, nil
}

func (s *rentService) Update(ctx context.Context, id uint64, lat, long float64, address, info string) (*model.Rent, error) {
	oldRent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, rent.ErrRentNotFound) {
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		return nil, err
	}
	if oldRent.Geopoint.Lat == lat && oldRent.Geopoint.Long == long && oldRent.Address == address && oldRent.Info == info {
		return oldRent, nil
	}
	if lat != 0 {
		err = oldRent.Geopoint.UpdateLat(lat)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLat) {
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			return nil, err
		}
	}
	if long != 0 {
		err = oldRent.Geopoint.UpdateLong(long)
		if err != nil {
			if errors.Is(err, apperrors.ErrInvalidLong) {
				return nil, apperrors.Wrap(err, apperrors.ValidationError(err.Error()))
			}
			return nil, err
		}
	}

	if address != "" {
		oldRent.UpdateAddress(address)
	}

	if info != "" {
		oldRent.UpdateInfo(info)
	}

	updatedRent, err := s.repo.Update(ctx, oldRent)
	if err != nil {
		if errors.Is(err, rent.ErrRentAlreadyExists) {
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		return nil, err
	}
	return updatedRent, nil
}

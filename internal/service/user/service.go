package user

import (
	"context"
	"errors"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type UserService struct {
	repo UserRepository
	log  logger.Logger
}

func NewUserService(log logger.Logger, repo UserRepository) *UserService {
	return &UserService{log: log, repo: repo}
}

func (s *UserService) Create(ctx context.Context, tgID uint64) (*model.User, error) {
	newUser := &model.User{TgID: tgID}

	newUser, err := s.repo.Create(ctx, newUser)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			s.log.Warn("user already exists", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrConflict)
		}
		s.log.Error("error creating user", "error", err)
		return nil, err
	}
	return newUser, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			s.log.Warn("user not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error getting user", "error", err)
		return nil, err
	}
	return result, nil
}

func (s *UserService) Update(ctx context.Context, id, newTgID uint64) (*model.User, error) {
	oldUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			s.log.Warn("user not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error updating user", "error", err)
		return nil, err
	}
	if oldUser.TgID != newTgID {
		oldUser.TgID = newTgID
		oldUser, err = s.repo.Update(ctx, oldUser)
		if err != nil {
			if errors.Is(err, user.ErrUserAlreadyExists) {
				s.log.Warn("user already exists", "error", err)
				return nil, apperrors.Wrap(err, apperrors.ErrConflict)
			}
			s.log.Error("error updating user", "error", err)
			return nil, err
		}
	}
	return oldUser, err
}

func (s *UserService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			s.log.Warn("user not found", "error", err)
			return apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error deleting user", "error", err)
		return err
	}
	return nil
}

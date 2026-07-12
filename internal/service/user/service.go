package user

import (
	"context"
	"errors"

	apperrors "github.com/kotafan1rich/geo_logic_api_go/internal/errors"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/user"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type userService struct {
	repo repository.UserRepository
	log  Logger
}

func NewUserService(log Logger, repo repository.UserRepository) service.UserService {
	return &userService{log: log, repo: repo}
}

func (s *userService) Create(ctx context.Context, tgId uint64) (*model.User, error) {
	newUser := &model.User{TgID: tgId}

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

func (s *userService) GetById(ctx context.Context, id uint64) (*model.User, error) {
	result, err := s.repo.GetById(ctx, id)
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

func (s *userService) Update(ctx context.Context, id, newTgId uint64) (*model.User, error) {
	oldUser, err := s.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			s.log.Warn("user not found", "error", err)
			return nil, apperrors.Wrap(err, apperrors.ErrNotFound)
		}
		s.log.Error("error updating user", "error", err)
		return nil, err
	}
	s.log.Info("user found", "user", oldUser)
	if oldUser.TgID != newTgId {
		oldUser.TgID = newTgId
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

func (s *userService) Delete(ctx context.Context, id uint64) error {
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

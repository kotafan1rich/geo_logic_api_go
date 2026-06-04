package user

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type Logger interface {
	Error(msg string, args ...any)
}

type userService struct {
	repo repository.UserRepository
	log  Logger
}

func NewUserService(repo repository.UserRepository) service.UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, tgId uint64) (*model.User, error) {
	newUser := &model.User{TgId: tgId}

	newUser, err := s.repo.Create(ctx, newUser)
	if err != nil {
		s.log.Error("error while creating new user", "err", err)
		return nil, err
	}
	return newUser, nil
}

func (s *userService) GetById(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.log.Error("error while getting user by id", "err", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint64, newTgId uint64) (*model.User, error) {
	oldUser, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if oldUser.TgId != newTgId {
		oldUser.TgId = newTgId
		oldUser, err = s.repo.Update(ctx, oldUser)
		if err != nil {
			s.log.Error("error while updating user", "err", err)
			return nil, err
		}
	}
	return oldUser, err
}

func (s *userService) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Error("error while deleting user", "err", err)
		return err
	}
	return nil
}

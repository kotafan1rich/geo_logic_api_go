package service

import (
	"context"
	"log/slog"

	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetById(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type UserService interface {
	Create(ctx context.Context, tgId uint64) (*model.User, error)
	GetById(ctx context.Context, id uint64) (*model.User, error)
	Update(ctx context.Context, id uint64, newTgId uint64) (*model.User, error)
	Delete(ctx context.Context, id uint64) error
}

type service struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, tgId uint64) (*model.User, error) {
	newUser := &model.User{TgId: tgId}

	newUser, err := s.repo.Create(ctx, newUser)
	if err != nil {
		slog.Error("error while creating new user", "err", err)
		return nil, err
	}
	return newUser, nil
}

func (s *service) GetById(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		slog.Error("error while getting user by id", "err", err)
		return nil, err
	}

	return user, nil
}

func (s *service) Update(ctx context.Context, id uint64, newTgId uint64) (*model.User, error) {
	oldUser, err := s.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if oldUser.TgId != newTgId {
		oldUser.TgId = newTgId
		oldUser, err = s.repo.Update(ctx, oldUser)
		if err != nil {
			slog.Error("error while updating user", "err", err)
			return nil, err
		}
	}
	return oldUser, err
}

func (s *service) Delete(ctx context.Context, id uint64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		slog.Error("error while deleting user", "err", err)
		return err
	}
	return nil
}

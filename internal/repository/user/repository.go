package repository

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
)

type userRepository struct {
	db database.DB
}

func NewRepository(db database.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.GORM().WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.GORM().WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.GORM().WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.GORM().WithContext(ctx).Delete(&model.User{}, id).Error
}

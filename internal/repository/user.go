package repository

import (
	"context"

	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/service"
)

type repository struct {
	db database.DB
}

func NewRepository(db database.DB) service.UserRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.GORM().WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) GetById(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.GORM().WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.GORM().WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *repository) Delete(ctx context.Context, id uint64) error {
	return r.db.GORM().WithContext(ctx).Delete(&model.User{}, id).Error
}

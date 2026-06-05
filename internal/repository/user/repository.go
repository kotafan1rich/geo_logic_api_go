package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository"
	"gorm.io/gorm"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetById(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := r.db.GORM().WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	result := r.db.GORM().WithContext(ctx).Model(user).Updates(user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrUserAlreadyExists
		}
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.GORM().WithContext(ctx).Delete(&model.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

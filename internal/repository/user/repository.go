package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kotafan1rich/geo_logic_api_go/internal/database"
	"github.com/kotafan1rich/geo_logic_api_go/internal/model"
	"github.com/kotafan1rich/geo_logic_api_go/internal/repository/user/dbmodel"
	"gorm.io/gorm"
)

type UserRepository struct {
	db database.DB
}

func NewRepository(db database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	userModel := dbmodel.ToUserModel(user)
	err := r.db.GORM().WithContext(ctx).Create(&userModel).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == database.ErrPgUniqueViolation {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	user.ID = userModel.ID
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user dbmodel.User
	err := r.db.GORM().WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dbmodel.ToUser(&user), nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	userModel := dbmodel.ToUserModel(user)
	result := r.db.GORM().WithContext(ctx).Model(&userModel).Select("*").Updates(&userModel)
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

func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	result := r.db.GORM().WithContext(ctx).Delete(&dbmodel.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

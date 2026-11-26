package repository

import (
	"context"
	"errors"

	"github.com/TwiLightDM/lab-track/internal/entity"
	errs "github.com/TwiLightDM/lab-track/internal/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	res := r.
		db.
		WithContext(ctx).
		Create(user)

	if res.Error != nil {
		switch {
		case errors.Is(res.Error, gorm.ErrDuplicatedKey):
			return errs.ErrDuplicatedKey
		default:
			return res.Error
		}
	}

	return nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	if err := r.
		db.
		WithContext(ctx).
		Find(&users).
		Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User

	if err := r.
		db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrEntityNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) PatchUser(ctx context.Context, id int, updates map[string]any) (*entity.User, error) {
	user := &entity.User{}

	res := r.
		db.
		WithContext(ctx).
		Model(user).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updates)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, errs.ErrEntityNotFound
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	res := r.
		db.
		WithContext(ctx).
		Delete(&entity.User{}, "id = ?", id)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errs.ErrEntityNotFound
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	if err := r.
		db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrEntityNotFound
		}

		return nil, err
	}

	return &user, nil
}

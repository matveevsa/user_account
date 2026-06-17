package repository

import (
	"account/internal/model"
	"account/internal/repository/mapper"
	repomodel "account/internal/repository/model"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
	lg *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db: db,
		lg: logger,
	}
}
func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	userRepo := mapper.UserToRepoUser(user)
	res := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: false}).
		Create(&userRepo)
	if res.Error != nil {
		r.lg.Error().Msgf("failed to save user")
		return fmt.Errorf("failed to save user: %w", res.Error)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	var user repomodel.User
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ?", userID).
		First(&user)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, fmt.Errorf("user not found")
	}

	if res.Error != nil {
		r.lg.Error().Msgf("failed to get user id")
		return model.User{}, fmt.Errorf("failed to get user id = %v: %w", userID, res.Error)
	}

	return mapper.RepoUserToUser(user), nil
}

func (r *Repository) GetUsers(ctx context.Context, offset, limit int) ([]model.User, error) {
	var users []repomodel.User
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Offset(offset).
		Limit(limit).
		Find(&users)

	if res.Error != nil {
		r.lg.Error().Msgf("failed to get users")
		return nil, fmt.Errorf("failed to get users: %w", res.Error)
	}

	return mapper.RepoUsersToUsers(users), nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uint64) error {
	res := r.db.WithContext(ctx).
		Where("id = ?", userID).
		Delete(&repomodel.User{})

	if res.Error != nil {
		r.lg.Error().Msgf("failed to delete user")
		return fmt.Errorf("failed to delete user id = %v: %w", userID, res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error {
	updates := map[string]interface{}{
		"email":       user.Email,
		"phone":       user.Phone,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"middle_name": user.MiddleName,
		"age":         user.Age,
		"updated_at":  user.UpdatedAt,
	}

	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ?", userID).
		Updates(updates)

	if res.Error != nil {
		r.lg.Error().Msgf("failed to update user")
		return fmt.Errorf("failed to update user: %w", res.Error)
	}

	return nil
}

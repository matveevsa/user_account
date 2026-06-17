package service

import (
	"account/internal/model"
	"context"
	"github.com/rs/zerolog"
	"time"
)

type AccountService struct {
	repo Repository
	lg   *zerolog.Logger
}

func New(repo Repository, lg *zerolog.Logger) *AccountService {
	return &AccountService{repo, lg}
}

type Repository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, userID uint64) (model.User, error)
	GetUsers(ctx context.Context, offset, limit int) ([]model.User, error)
	DeleteUser(ctx context.Context, userID uint64) error
	UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error
}

func (as *AccountService) CreateUser(ctx context.Context, newUser model.CreateUser) error {
	user := model.User{
		Login:      newUser.Login,
		Email:      newUser.Email,
		Phone:      newUser.Phone,
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		MiddleName: newUser.MiddleName,
		Age:        newUser.Age,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return as.repo.CreateUser(ctx, user)
}

func (as *AccountService) GetUsers(ctx context.Context, offset, limit int) ([]model.User, error) {
	return as.repo.GetUsers(ctx, offset, limit)
}

func (as *AccountService) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	return as.repo.GetUser(ctx, userID)
}

func (as *AccountService) DeleteUser(ctx context.Context, userID uint64) error {
	return as.repo.DeleteUser(ctx, userID)
}

func (as *AccountService) UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error {
	user.UpdatedAt = time.Now()
	return as.repo.UpdateUser(ctx, userID, user)
}

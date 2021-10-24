package service

import (
	"context"
	"taskESchop/internal/models"
	"taskESchop/internal/repository"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (u UserServiceImpl) AddProductToCart(ctx context.Context, req *models.Product, userId int) error {
	return u.repo.AddProductToCart(ctx, req, userId)
}

func (u UserServiceImpl) DeleteFromCart(ctx context.Context, id int) error {
	return u.repo.DeleteFromCart(ctx, id)
}

func (u UserServiceImpl) GetCart(ctx context.Context, userId int) ([]*models.GetCartRes, error) {
	return u.repo.GetCart(ctx, userId)
}

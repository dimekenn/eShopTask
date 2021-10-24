package service

import (
	"context"
	"taskESchop/internal/models"
)

type UserService interface {
	AddProductToCart(ctx context.Context, req *models.Product, userId int) error
	DeleteFromCart(ctx context.Context, id int) error
	GetCart(ctx context.Context, userId int) ([]*models.GetCartRes, error)
}

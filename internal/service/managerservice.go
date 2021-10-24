package service

import (
	"context"
	"taskESchop/internal/models"
)

type ManagerService interface {
	CreateNewProduct(ctx context.Context, req *models.Product) (int, error)
	UpdateProduct(ctx context.Context, req *models.Product) error
	GetAllCarts(ctx context.Context) ([]*models.Cart, error)
}

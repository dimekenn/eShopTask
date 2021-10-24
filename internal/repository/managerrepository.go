package repository

import (
	"golang.org/x/net/context"
	"taskESchop/internal/models"
)

type ManagerRepository interface {
	CreateNewProduct(ctx context.Context, req *models.Product) (int, error)
	UpdateProduct(ctx context.Context, req *models.Product) error
	GetAllCarts(ctx context.Context) ([]*models.DBCart, error)
}

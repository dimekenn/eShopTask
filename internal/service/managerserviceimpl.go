package service

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"strings"
	"taskESchop/internal/models"
	"taskESchop/internal/repository"
)

type ManagerServiceImpl struct {
	repo repository.ManagerRepository
}

func NewManagerService(repo repository.ManagerRepository) ManagerService  {
	return &ManagerServiceImpl{repo: repo}
}

func (m ManagerServiceImpl) CreateNewProduct(ctx context.Context, req *models.Product) (int, error) {
	return m.repo.CreateNewProduct(ctx, req)
}

func (m ManagerServiceImpl) UpdateProduct(ctx context.Context, req *models.Product) error {
	return m.repo.UpdateProduct(ctx, req)
}

func (m ManagerServiceImpl) GetAllCarts(ctx context.Context) ([]*models.Cart, error) {
	carts, err := m.repo.GetAllCarts(ctx)
	if err != nil{
		return nil, err
	}
	cartArr := make([]*models.Cart, len(carts))
	for i, v  := range carts{
		cartArr[i] = &models.Cart{}
		cartArr[i].Id = v.Id
		cartArr[i].UserId = v.UserId
		cartArr[i].Products = make([]*models.CartProduct, len(v.Products))
		for j, k := range v.Products{
			cartArr[i].Products[j] = &models.CartProduct{}
			if k != ""{
				pArr := strings.Split(k, ":")
				id, idErr := strconv.Atoi(pArr[0])
				if idErr!=nil{
					log.Errorf("parse cart product's id error: %v", idErr)
					return nil, echo.NewHTTPError(http.StatusInternalServerError, idErr)
				}
				productId, pIdErr := strconv.Atoi(pArr[1])
				if pIdErr!=nil{
					log.Errorf("parse product's id error: %v", pIdErr)
					return nil, echo.NewHTTPError(http.StatusInternalServerError, pIdErr)
				}
				count, cErr := strconv.Atoi(pArr[2])
				if cErr!=nil{
					log.Errorf("parse cart product's count error: %v", cErr)
					return nil, echo.NewHTTPError(http.StatusInternalServerError, cErr)
				}
				cartArr[i].Products[j].Id = id
				cartArr[i].Products[j].ProductId = productId
				cartArr[i].Products[j].Count = count
			}
		}
	}
	return cartArr, nil
}

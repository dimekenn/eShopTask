package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"taskESchop/internal/models"
)

type UserRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository  {
	return &UserRepositoryImpl{pool: pool}
}

func (u UserRepositoryImpl) AddProductToCart(ctx context.Context, req *models.Product, userId int) error {
	tx, txErr := u.pool.Begin(ctx)
	if txErr!=nil{
		log.Errorf("Error create tx in AddProductToCart: %v", txErr)
		return echo.NewHTTPError(http.StatusInternalServerError, txErr)
	}

	_, err := tx.Exec(
		ctx,
		"insert into cart_products (cart_id, product_id, count) values ((select id from carts where user_id = $1), $2, $3) returning id",
		userId, req.Id, req.Count,
		)
	if err != nil{
		rbErr := tx.Rollback(ctx)
		if rbErr != nil{
			log.Errorf("Error rollback tx in AddProductToCart: %v", txErr)
			return echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("Error exec in AddProductToCart: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	cErr := tx.Commit(ctx)
	if cErr != nil{
		log.Errorf("Error commit tx in AddProductToCart: %v", txErr)
		return echo.NewHTTPError(http.StatusInternalServerError, cErr)
	}
	return nil
}

func (u UserRepositoryImpl) DeleteFromCart(ctx context.Context, id int) error {
	tx, txErr := u.pool.Begin(ctx)
	if txErr!=nil{
		log.Errorf("Error create tx in DeleteFromCart: %v", txErr)
		return echo.NewHTTPError(http.StatusInternalServerError, txErr)
	}

	_, err := tx.Exec(
		ctx,
		"delete from cart_products where id = $1",
		id,
		)
	if err != nil{
		rbErr := tx.Rollback(ctx)
		if rbErr != nil{
			log.Errorf("Error rollback tx in DeleteFromCart: %v", txErr)
			return echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("Error exec in DeleteFromCart: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	cErr := tx.Commit(ctx)
	if cErr != nil{
		log.Errorf("Error commit tx in AddProductToCart: %v", txErr)
		return echo.NewHTTPError(http.StatusInternalServerError, cErr)
	}
	return nil
}

func (u UserRepositoryImpl) GetCart(ctx context.Context, userId int) ([]*models.GetCartRes, error) {
	var res []*models.GetCartRes
	rows, err := u.pool.Query(
		ctx,
		"select cp.id, cp.product_id, cp.cart_id, cp.count from cart_products cp join carts c on c.id = cp.cart_id where c.user_id = $1",
		userId,
		)
	if err != nil{
		if err == pgx.ErrNoRows{
			return nil, echo.NewHTTPError(http.StatusOK, "cart is empty")
		}
		log.Errorf("Error exec in GetCart: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for rows.Next(){
		cartProduct := &models.GetCartRes{}
		scErr := rows.Scan(&cartProduct.Id, &cartProduct.ProductId, &cartProduct.CartId, &cartProduct.Count)
		if scErr != nil{
			log.Errorf("Error scan object in GetCart: %v", scErr)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, scErr)
		}
		res = append(res, cartProduct)
	}
	return res, nil
}

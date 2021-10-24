package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"net/http"
	"taskESchop/internal/models"
)

type ManagerRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewManagerRepository(pool *pgxpool.Pool) ManagerRepository  {
	return &ManagerRepositoryImpl{pool: pool}
}

func (m ManagerRepositoryImpl) CreateNewProduct(ctx context.Context, req *models.Product) (int, error) {
	tx, txErr := m.pool.Begin(ctx)
	if txErr!=nil{
		log.Errorf("Error create tx in CreateNewProduct: %v", txErr)
		return 0, txErr
	}

	var id int

	err := tx.QueryRow(
		ctx,
		"insert into products (name, description, count) values ($1, $2, $3) returning id",
		req.Name, req.Description,req.Count,
		).Scan(&id)
	if err != nil{
		rbErr := tx.Rollback(ctx)
		if rbErr != nil{
			log.Errorf("Error rollback tx in CreateNewProduct: %v", txErr)
			return 0, echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("Error exec in CreateNewProduct: %v", err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	cErr := tx.Commit(ctx)
	if cErr!=nil{
		log.Errorf("Error commit tx in CreateNewProduct: %v", cErr)
		return 0, echo.NewHTTPError(http.StatusInternalServerError, cErr)
	}
	return id, nil
}

func (m ManagerRepositoryImpl) UpdateProduct(ctx context.Context, req *models.Product)  error {
	tx, txErr := m.pool.Begin(ctx)
	if txErr!=nil{
		log.Errorf("Error create tx in UpdateProduct: %v", txErr)
		return txErr
	}

	_, err := tx.Exec(
		ctx,
		"update products set name = coalesce(TRIM($1), name), description = coalesce(TRIM($2), description), count = $3 where id = $4",
		req.Name, req.Description, req.Count, req.Id,
		)

	if err != nil{
		rbErr := tx.Rollback(ctx)
		if rbErr != nil{
			log.Errorf("Error rollback tx in UpdateProduct: %v", txErr)
			return echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("Error exec in UpdateProduct: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	cErr := tx.Commit(ctx)
	if cErr!=nil{
		log.Errorf("Error commit tx in CreateNewProduct: %v", cErr)
		return echo.NewHTTPError(http.StatusInternalServerError, cErr)
	}
	return nil
}

func (m ManagerRepositoryImpl) GetAllCarts(ctx context.Context) ([]*models.DBCart, error) {
	rows, err := m.pool.Query(
		ctx,
		"select c.id, c.user_id, array_agg(cp.id || ':' || cp.product_id || ':' || cp.count) from carts c left outer join cart_products cp on cp.cart_id = c.id group by c.id, c.user_id",
		)
	if err != nil{
		if err == pgx.ErrNoRows{
			return nil, echo.NewHTTPError(http.StatusOK, "not found")
		}
		log.Errorf("Error query in GetAllCarts: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var carts []*models.DBCart
	var nilArr []sql.NullString

	for rows.Next(){
		cart := &models.DBCart{}
		scErr := rows.Scan(&cart.Id, &cart.UserId, pq.Array(&nilArr))
		if scErr != nil{
			log.Errorf("Error scan objects: %v", scErr)
			return nil, echo.NewHTTPError(http.StatusInternalServerError, scErr)
		}
		cart.Products = make([]string, len(nilArr))
		for i, v := range nilArr{
			cart.Products[i] = v.String
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

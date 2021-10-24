package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"taskESchop/internal/models"
)

type AuthorizationRepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewAuthorizationRepository(pool *pgxpool.Pool) AuthorizationRepository {
	return &AuthorizationRepositoryImpl{pool: pool}
}

func (a AuthorizationRepositoryImpl) SignUp(ctx context.Context, user *models.User, password []byte) error {
	tx, txErr := a.pool.Begin(ctx)
	if txErr!=nil{
		log.Errorf("Error create tx in SignUp: %v", txErr)
		return echo.NewHTTPError(http.StatusInternalServerError, txErr)
	}

	buff := bytes.NewBuffer(make([]byte, 0))

	for i, v := range user.Roles{
		buff.WriteString(strconv.Itoa(v))
		if i != len(user.Roles)-1{
			buff.WriteString(",")
		}
	}

	query := fmt.Sprintf("with new_user as (insert into users(username, password) values($1, $2) returning id) insert into user_roles values((select id from new_user), unnest(array[%v]))", buff.String())

	_, err := tx.Exec(
		ctx,
		query,
		user.Username, password,
		)
	if err != nil{
		rbErr := tx.Rollback(ctx)
		if rbErr != nil{
			log.Errorf("Error rollback tx in SignUp: %v", txErr)
			return echo.NewHTTPError(http.StatusInternalServerError, rbErr)
		}
		log.Errorf("Error exec in SignUp: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	cErr := tx.Commit(ctx)
	if cErr!=nil{
		log.Errorf("Error commit tx in SignUp: %v", cErr)
		return echo.NewHTTPError(http.StatusInternalServerError, cErr)
	}
	return nil
}

func (a AuthorizationRepositoryImpl) SignIn(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := a.pool.QueryRow(
		ctx,
		"select u.id, u.username, u.password, array_agg(r.name) from users u join user_roles ur on ur.user_id = u.id join roles r on r.id = ur.role_id where u.username = $1 group by u.id, u.username , u.password",
		username,
		).Scan(&user.Id, &user.Username, &user.Password, pq.Array(&user.RolesName))

	if err != nil{
		if err == pgx.ErrNoRows{
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "username not found")
		}
		log.Errorf("Error query row in SignIn: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return user, nil
}

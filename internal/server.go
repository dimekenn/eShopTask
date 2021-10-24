package internal

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"taskESchop/internal/configs"
	"taskESchop/internal/handler"
	"taskESchop/internal/middlewares"
	"taskESchop/internal/repository"
	"taskESchop/internal/service"
	"time"
)

func StartHTTPServer(ctx context.Context, errCh chan<- error, cfg *configs.Configs) {
	app := echo.New()

	pool, poolErr := initDb(ctx, cfg.DbUrl)
	if poolErr != nil {
		errCh <- poolErr
		return
	}

	mw := middlewares.NewAuthorizationMiddleware()
	app.Use(mw.Authorize)

	authorizationRepo := repository.NewAuthorizationRepository(pool)
	managerRepo := repository.NewManagerRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	authorizationService := service.NewAuthorizationService(authorizationRepo)
	managerService := service.NewManagerService(managerRepo)
	userService := service.NewUserService(userRepo)

	srvHandler := handler.NewHandler(userService, managerService, authorizationService)

	app.POST("/signup", srvHandler.SignUp)
	app.POST("/signin", srvHandler.SignIn)
	app.POST("/product", srvHandler.CreateNewProduct)
	app.PUT("/product", srvHandler.UpdateProduct)
	app.GET("/carts", srvHandler.GetAllCarts)
	app.GET("/cart", srvHandler.GetCart)
	app.DELETE("/cart/:id", srvHandler.DeleteFromCart)
	app.POST("/cart", srvHandler.AddProductToCart)

	errCh <- app.Start(cfg.Port)
}

func initDb(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conf, cfgErr := pgxpool.ParseConfig(url)
	if cfgErr != nil {
		return nil, cfgErr
	}
	conf.MaxConns = 20
	conf.MinConns = 10
	conf.MaxConnIdleTime = 10 * time.Second

	conn, connErr := pgxpool.ConnectConfig(ctx, conf)
	if connErr != nil {
		return nil, connErr
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

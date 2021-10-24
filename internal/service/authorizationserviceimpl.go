package service

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"taskESchop/internal/models"
	"taskESchop/internal/repository"
)

type AuthorizationServiceImpl struct {
	repo repository.AuthorizationRepository
}

func NewAuthorizationService(repo repository.AuthorizationRepository) AuthorizationService {
	return &AuthorizationServiceImpl{repo: repo}
}

func (a AuthorizationServiceImpl) SignUp(ctx context.Context, req *models.User) (*models.Msg, error) {
	pwBytes, bcErr := bcrypt.GenerateFromPassword([]byte(req.Password), 5)
	if bcErr!=nil{
		log.Errorf("Error GenerateFromPassword: %v", bcErr)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, bcErr)
	}
	err := a.repo.SignUp(ctx, req, pwBytes)
	if err != nil{
		return nil, err
	}
	return &models.Msg{Msg: "success created"}, nil

}

func (a AuthorizationServiceImpl) SignIn(ctx context.Context, req *models.User) (*models.User, error) {
	user, err := a.repo.SignIn(ctx, req.Username)
	if err != nil{
		return nil, err
	}
	bcErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if bcErr != nil{
		log.Warnf("Password error: %v", bcErr)
		return nil, echo.NewHTTPError(http.StatusUnauthorized, bcErr)
	}
	user.Password = ""
	return user, nil
}

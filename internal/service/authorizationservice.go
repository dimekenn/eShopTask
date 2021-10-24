package service

import (
	"context"
	"taskESchop/internal/models"
)

type AuthorizationService interface {
	SignUp(ctx context.Context, user *models.User) (*models.Msg, error)
	SignIn(ctx context.Context, user *models.User) (*models.User, error)
}

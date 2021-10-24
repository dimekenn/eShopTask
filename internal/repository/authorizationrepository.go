package repository

import (
	"context"
	"taskESchop/internal/models"
)

type AuthorizationRepository interface {
	SignUp(ctx context.Context, user *models.User, password []byte) error
	SignIn(ctx context.Context, username string) (*models.User, error)
}

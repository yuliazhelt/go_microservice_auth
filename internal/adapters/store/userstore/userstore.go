package userstore

import (
	"context"
	"auth/internal/models"
)

type User interface {
	Get(ctx context.Context, login string) (*models.User, error)
	// Add(ctx context.Context, login string, password string, role string) (uint64, error)
}

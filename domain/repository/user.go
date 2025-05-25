package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/internal/shared"
)

type UserRepositorier interface {
	Get(ctx context.Context, profile string, op shared.QueryOptions) ([]model.User, error)

	Insert(ctx context.Context, profile string, users *[]model.User) error

	UpdateCount(ctx context.Context, profile string) error
}

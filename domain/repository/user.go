package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/internal/common"
)

type UserRepositorier interface {
	Get(ctx context.Context, profile string, op common.QueryOptions) ([]model.User, error)

	Insert(ctx context.Context, profile string, users *[]model.User) error

	UpdateCount(ctx context.Context, profile string) error
}

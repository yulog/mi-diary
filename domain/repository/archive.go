package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
)

type ArchiveRepositorier interface {
	Get(ctx context.Context, profile string) ([]model.Month, error)

	UpdateCountMonthly(ctx context.Context, profile string) error
	UpdateCountDaily(ctx context.Context, profile string) error
}

package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/util/pagination"
)

type FileRepositorier interface {
	Get(ctx context.Context, profile, c string, p pagination.Paging) ([]model.File, error)
	GetByNoteID(ctx context.Context, profile, id string) ([]model.File, error)
	GetByEmptyColor(ctx context.Context, profile string) ([]model.File, error)

	Count(ctx context.Context, profile string) (int, error)

	Insert(ctx context.Context, profile string, files *[]model.File) error

	UpdateByPKWithColor(ctx context.Context, profile, id, c1, c2 string)
}

package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
)

type HashTagRepositorier interface {
	Get(ctx context.Context, profile string) ([]model.HashTag, error)

	Insert(ctx context.Context, profile string, hashtag *model.HashTag) error
}

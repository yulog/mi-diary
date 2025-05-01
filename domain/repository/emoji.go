package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	mi "github.com/yulog/miutil"
)

type EmojiRepositorier interface {
	Get(ctx context.Context, profile string) ([]model.ReactionEmoji, error)
	GetByName(ctx context.Context, profile, name string) (model.ReactionEmoji, error)
	GetByEmptyImage(ctx context.Context, profile string) ([]model.ReactionEmoji, error)

	Insert(ctx context.Context, profile string, reactions *[]model.ReactionEmoji) error

	UpdateByPKWithImage(ctx context.Context, profile string, id int64, e *mi.Emoji)
	UpdateCount(ctx context.Context, profile string) error
}

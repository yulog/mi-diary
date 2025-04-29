package infra

import (
	"context"
	"log/slog"

	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/model"
	mi "github.com/yulog/miutil"
)

type EmojiInfra struct {
	infra *Infra
}

func (i *Infra) NewEmojiInfra() logic.EmojiRepositorier {
	return &EmojiInfra{infra: i}
}

func (ei *EmojiInfra) Get(ctx context.Context, profile string) ([]model.ReactionEmoji, error) {
	var reactions []model.ReactionEmoji
	err := ei.infra.DB(profile).
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (ei *EmojiInfra) GetByName(ctx context.Context, profile, name string) (model.ReactionEmoji, error) {
	var reaction model.ReactionEmoji
	err := ei.infra.DB(profile).
		NewSelect().
		Model(&reaction).
		Where("name = ?", name).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return model.ReactionEmoji{}, err
	}
	return reaction, nil
}

func (ei *EmojiInfra) GetByEmptyImage(ctx context.Context, profile string) ([]model.ReactionEmoji, error) {
	var reactions []model.ReactionEmoji
	err := ei.infra.DB(profile).
		NewSelect().
		Model(&reactions).
		Where("image = ?", "").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (ei *EmojiInfra) Insert(ctx context.Context, profile string, reactions *[]model.ReactionEmoji) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = ei.infra.DB(profile)
	}
	_, err := db.NewInsert().
		Model(reactions).
		Ignore().
		Exec(ctx)
	return err
}

func (ei *EmojiInfra) UpdateByPKWithImage(ctx context.Context, profile string, id int64, e *mi.Emoji) {
	// TODO: emoji画像をローカルに保存する

	r := model.ReactionEmoji{
		ID:    id,
		Image: e.URL,
	}
	var s []model.ReactionEmoji
	s = append(s, r)
	_, err := ei.infra.DB(profile).NewUpdate().
		Model(&s).
		OmitZero().
		Column("image").
		Bulk().
		Exec(ctx)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

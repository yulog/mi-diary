package infra

import (
	"context"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
	mi "github.com/yulog/miutil"
)

type EmojiInfra struct {
	dao *DataBase
}

func (i *Infra) NewEmojiInfra() repository.EmojiRepositorier {
	return &EmojiInfra{dao: i.dao}
}

func (i *EmojiInfra) Get(ctx context.Context, profile string) ([]model.ReactionEmoji, error) {
	var reactions []model.ReactionEmoji
	err := i.dao.DB(profile).
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (i *EmojiInfra) GetByName(ctx context.Context, profile, name string) (model.ReactionEmoji, error) {
	var reaction model.ReactionEmoji
	err := i.dao.DB(profile).
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

func (i *EmojiInfra) GetByEmptyImage(ctx context.Context, profile string) ([]model.ReactionEmoji, error) {
	var reactions []model.ReactionEmoji
	err := i.dao.DB(profile).
		NewSelect().
		Model(&reactions).
		Where("image = ?", "").
		Where("is_symbol = ?", false).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (i *EmojiInfra) Insert(ctx context.Context, profile string, reactions *[]model.ReactionEmoji) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.dao.DB(profile)
	}
	_, err := db.NewInsert().
		Model(reactions).
		Ignore().
		Exec(ctx)
	return err
}

func (i *EmojiInfra) UpdateByPKWithImage(ctx context.Context, profile string, id int64, e *mi.Emoji) {
	// TODO: emoji画像をローカルに保存する

	r := model.ReactionEmoji{
		ID:    id,
		Image: e.URL,
	}
	var s []model.ReactionEmoji
	s = append(s, r)
	_, err := i.dao.DB(profile).NewUpdate().
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

// リアクションのカウント
func (i *EmojiInfra) UpdateCount(ctx context.Context, profile string) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.dao.DB(profile)
	}
	var reactions []model.ReactionEmoji
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		Relation("Reaction", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ColumnExpr("reaction.id as id")
		}).
		ColumnExpr("reaction_emoji_name as name").
		ColumnExpr("count(*) as count").
		Group("reaction_emoji_name").
		Scan(ctx, &reactions)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().
		Model(&reactions).
		OmitZero().
		Column("count").
		Bulk().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

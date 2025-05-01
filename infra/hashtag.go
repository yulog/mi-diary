package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
)

type HashTagInfra struct {
	infra *DataBase
}

func (i *Infra) NewHashTagInfra() repository.HashTagRepositorier {
	return &HashTagInfra{infra: i.dao}
}

func (hi *HashTagInfra) Get(ctx context.Context, profile string) ([]model.HashTag, error) {
	var tags []model.HashTag
	err := hi.infra.DB(profile).
		NewSelect().
		Model(&tags).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (hi *HashTagInfra) Insert(ctx context.Context, profile string, hashtag *model.HashTag) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = hi.infra.DB(profile)
	}
	_, err := db.NewInsert().
		Model(hashtag).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	return err
}

// タグのカウント
func (hi *HashTagInfra) UpdateCount(ctx context.Context, profile string) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = hi.infra.DB(profile)
	}
	var hashtags []model.HashTag
	err := db.NewSelect().
		Model((*model.NoteToTag)(nil)).
		Relation("HashTag", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("text")
		}).
		ColumnExpr("hash_tag_id as id").
		ColumnExpr("count(*) as count").
		Group("hash_tag_id").
		Scan(ctx, &hashtags)
	if err != nil {
		return err
	}

	if len(hashtags) > 0 {
		_, err = db.NewUpdate().
			Model(&hashtags).
			OmitZero().
			Column("count").
			Bulk().
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

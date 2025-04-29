package infra

import (
	"context"

	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/model"
)

type HashTagInfra struct {
	infra *Infra
}

func (i *Infra) NewHashTagInfra() logic.HashTagRepositorier {
	return &HashTagInfra{infra: i}
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

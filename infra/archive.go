package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
)

type ArchiveInfra struct {
	infra *Infra
}

func (i *Infra) NewArchiveInfra() repository.ArchiveRepositorier {
	return &ArchiveInfra{infra: i}
}

func (i *ArchiveInfra) Get(ctx context.Context, profile string) ([]model.Month, error) {
	var archives []model.Month
	err := i.infra.DB(profile).
		NewSelect().
		Model(&archives).
		Relation("Days", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ymd DESC")
		}).
		Order("ym DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return archives, nil
}

// 月別のカウント
func (i *ArchiveInfra) UpdateCountMonthly(ctx context.Context, profile string) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.infra.DB(profile)
	}
	var months []model.Month
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym").
		ColumnExpr("count(*) as count").
		Group("ym").
		Having("ym is not null").
		Scan(ctx, &months)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().
		Model(&months).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 日別のカウント
func (i *ArchiveInfra) UpdateCountDaily(ctx context.Context, profile string) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.infra.DB(profile)
	}
	var days []model.Day
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("strftime('%Y-%m-%d', created_at, 'localtime') as ymd").
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym").
		ColumnExpr("count(*) as count").
		Group("ymd").
		Having("ymd is not null").
		Scan(ctx, &days)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().
		Model(&days).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

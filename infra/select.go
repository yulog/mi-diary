package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
)

// addWhereLike
//
// https://bun.uptrace.dev/guide/query-where.html#querybuilder
func addWhereLike(q bun.QueryBuilder, col, s string) bun.QueryBuilder {
	if s == "" {
		return q
	}
	return q.Where("? like ?", bun.Ident(col), "%"+s+"%")
}

func addWhere(q bun.QueryBuilder, col, s string) bun.QueryBuilder {
	if s == "" {
		return q
	}
	return q.Where("? = ?", bun.Ident(col), s)
}

func (infra *Infra) Archives(ctx context.Context, profile string) ([]model.Month, error) {
	var archives []model.Month
	err := infra.DB(profile).
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

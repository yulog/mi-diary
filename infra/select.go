package infra

import (
	"github.com/uptrace/bun"
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

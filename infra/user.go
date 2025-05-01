package infra

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/internal/common"
)

type UserInfra struct {
	infra *Infra
}

func (i *Infra) NewUserInfra() repository.UserRepositorier {
	return &UserInfra{infra: i}
}

func (ui *UserInfra) Get(ctx context.Context, profile string, op common.QueryOptions) ([]model.User, error) {
	var users []model.User
	err := ui.infra.DB(profile).
		NewSelect().
		Model(&users).
		Order(fmt.Sprintf("%s %s", op.SortBy, op.SortOrder)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (i *UserInfra) Insert(ctx context.Context, profile string, users *[]model.User) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.infra.DB(profile)
	}
	_, err := db.NewInsert().
		Model(users).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	return err
}

// ユーザーのカウント
func (i *UserInfra) UpdateCount(ctx context.Context, profile string) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.infra.DB(profile)
	}
	var users []model.User
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name")
		}).
		ColumnExpr("count(*) as count").
		Group("user_id").
		Scan(ctx, &users)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().
		Model(&users).
		OmitZero().
		Column("count").
		Bulk().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

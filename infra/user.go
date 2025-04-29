package infra

import (
	"context"
	"fmt"

	"github.com/yulog/mi-diary/internal/common"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/model"
)

type UserInfra struct {
	infra *Infra
}

func (i *Infra) NewUserInfra() logic.UserRepositorier {
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

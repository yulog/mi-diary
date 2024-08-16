package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/model"
)

type UserInfra struct {
	infra *Infra
}

func (i *Infra) NewUserInfra() logic.UserRepositorier {
	return &UserInfra{infra: i}
}

func (ui *UserInfra) Get(ctx context.Context, profile string) ([]model.User, error) {
	var users []model.User
	err := ui.infra.DB(profile).
		NewSelect().
		Model(&users).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (infra *UserInfra) Insert(ctx context.Context, db bun.IDB, users *[]model.User) error {
	_, err := db.NewInsert().
		Model(users).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	return err
}

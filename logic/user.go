package logic

import (
	"context"

	"github.com/a-h/templ"
	"github.com/uptrace/bun"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/internal/common"
	"github.com/yulog/mi-diary/model"
)

type UserRepositorier interface {
	Get(ctx context.Context, profile string, op common.QueryOptions) ([]model.User, error)

	Insert(ctx context.Context, db bun.IDB, users *[]model.User) error
}

// type UserLogic struct {
// Repo UserRepositorier
// }

func (l *Logic) UsersLogic(ctx context.Context, profile string, partial bool, sortBy string) (templ.Component, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	var u []model.User
	if sortBy == "name" {
		u, err = l.UserRepo.Get(ctx, profile, common.QueryOptions{SortBy: "name", SortOrder: "ASC"})
	} else {
		u, err = l.UserRepo.Get(ctx, profile, common.QueryOptions{SortBy: "count", SortOrder: "DESC"})
	}
	if err != nil {
		return nil, err
	}
	return cm.Users(profile, u), nil
}

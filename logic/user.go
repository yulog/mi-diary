package logic

import (
	"context"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/internal/common"
)

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

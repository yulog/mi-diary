package logic

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/internal/shared"
)

func (l *Logic) UsersLogic(ctx context.Context, profile string, partial bool, sortBy string) (*UserOutput, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	var u []model.User
	if sortBy == "name" {
		u, err = l.UserRepo.Get(ctx, profile, shared.SortOptions{SortBy: "name", SortOrder: "ASC"})
	} else {
		u, err = l.UserRepo.Get(ctx, profile, shared.SortOptions{SortBy: "count", SortOrder: "DESC"})
	}
	if err != nil {
		return nil, err
	}
	return &UserOutput{
		Profile: profile,
		Users:   u,
	}, nil
}

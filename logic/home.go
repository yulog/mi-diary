package logic

import (
	"context"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/internal/common"
)

func (l *Logic) HomeLogic(ctx context.Context, profile string) (templ.Component, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	r, err := l.EmojiRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}

	return cm.IndexParams{
		Title:     profile,
		Profile:   profile,
		Reactions: r,
	}.Index(), nil
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile string) (templ.Component, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	h, err := l.HashTagRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}
	return cm.HashTags(profile, h), nil
}

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

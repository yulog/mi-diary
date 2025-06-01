package logic

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/internal/shared"
)

func (l *Logic) HomeLogic(ctx context.Context, profile string) (*IndexOutput, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	r, err := l.EmojiRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}

	return &IndexOutput{
		Title:     profile,
		Profile:   profile,
		Reactions: r,
	}, nil
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile string) (*HashTagOutput, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	h, err := l.HashTagRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}
	return &HashTagOutput{
		Profile:  profile,
		HashTags: h,
	}, nil
}

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

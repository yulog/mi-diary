package logic

import (
	"context"
)

func (l *Logic) ArchivesLogic(ctx context.Context, profile string) (*ArchivesOutput, error) {
	a, err := l.ArchiveRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}
	return &ArchivesOutput{
		Title:   "Archives",
		Profile: profile,
		Items:   a,
	}, nil
}

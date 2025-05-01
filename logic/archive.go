package logic

import (
	"context"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
)

func (l *Logic) ArchivesLogic(ctx context.Context, profile string) (templ.Component, error) {
	a, err := l.ArchiveRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}
	return cm.ArchiveParams{
		Title:   "Archives",
		Profile: profile,
		Items:   a,
	}.Archive(), nil
}

package logic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/yulog/mi-diary/internal/shared"
	"github.com/yulog/mi-diary/util/pagination"
)

func (l *Logic) FilesLogic(ctx context.Context, profile string, params shared.QueryParams) (*FileWithPages, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	count := 0
	if params.Color == "" {
		count, err = l.FileRepo.Count(ctx, profile)
		if err != nil {
			return nil, err
		}
		slog.Info("File count", slog.Int("count", count))
	}

	p, err := pagination.New(params.Page, 10, count, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("page count", slog.Int("count", p.CurrentPage))

	files, err := l.FileRepo.Get(ctx, profile, params.Color, p)
	if err != nil {
		return nil, err
	}
	slog.Info("file result count", slog.Int("count", len(files)))
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found")
	}

	if params.Color != "" {
		p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(files)}
	}
	slog.Info("has next", slog.Bool("bool", p.HasNextPage()))

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	hasLast := p.CurrentPage+1 < p.TotalPages()
	slog.Info("has last", slog.Bool("bool", hasLast))

	return &FileWithPages{
		File: File{
			Title:          fmt.Sprint(p.CurrentPage),
			Profile:        profile,
			Host:           host,
			FileFilterPath: fmt.Sprintf("/profiles/%s/files", profile),
			Items:          files,
		},
		Pages: Pages{
			Current: p.CurrentPage,
			Prev:    Page{Index: prev, Has: p.HasPreviousPage()},
			Next:    Page{Index: next, Has: p.HasNextPage()},
			Last:    Page{Index: p.TotalPages(), Has: hasLast},
			QueryParams: shared.QueryParams{
				Page:  params.Page,
				Color: params.Color,
			},
		},
	}, nil
}

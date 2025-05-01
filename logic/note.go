package logic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/util/pagination"
)

func (l *Logic) ReactionNotesLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.NoteRepo.GetByReaction(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) HashTagNotesLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.NoteRepo.GetByHashTag(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) UserLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.NoteRepo.GetByUser(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", name, p.CurrentPage),
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) NotesLogic(ctx context.Context, profile string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	count := 0
	if params.S == "" {
		count, err = l.NoteRepo.Count(ctx, profile)
		if err != nil {
			return nil, err
		}
	}

	p, err := pagination.New(params.Page, 10, count, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.NoteRepo.Get(ctx, profile, params.S, p)
	if err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, fmt.Errorf("note not found")
	}
	title := ""
	if params.S != "" {
		title = fmt.Sprintf("%s - %d", params.S, p.CurrentPage)
		p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}
	} else {
		title = fmt.Sprint(p.CurrentPage)
	}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	hasLast := p.CurrentPage+1 < p.TotalPages()
	// slog.Info("has last", slog.Bool("hasLast", hasLast), slog.Int("next", next), slog.Int("total", p.TotalPages()), slog.Int("current", p.CurrentPage))

	n := cm.Note{
		Title:      title,
		Profile:    profile,
		Host:       host,
		SearchPath: fmt.Sprintf("/profiles/%s/notes", profile),
		Items:      notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages(), Has: hasLast},
		QueryParams: cm.QueryParams{
			Page: params.Page,
			S:    params.S,
		},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) ArchiveNotesLogic(ctx context.Context, profile, d string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}
	// slog.Info("perPage", slog.Int("perPage", p2.Limit()))

	notes, err := l.NoteRepo.GetByArchive(ctx, profile, d, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", d, p.CurrentPage),
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
	}

	return n.WithPages(cp), nil
}

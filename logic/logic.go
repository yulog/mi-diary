package logic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/infra"
	"github.com/yulog/mi-diary/util/pg"
)

type Logic struct {
	repo *infra.Infra
}

func New(r *infra.Infra) *Logic {
	return &Logic{repo: r}
}

type Params struct {
	Page int
	S    string
}

func (l *Logic) HomeLogic(ctx context.Context, profile string) (templ.Component, error) {
	_, err := l.repo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	r, err := l.repo.Reactions(ctx, profile)
	if err != nil {
		return nil, err
	}
	h, err := l.repo.HashTags(ctx, profile)
	if err != nil {
		return nil, err
	}
	u, err := l.repo.Users(ctx, profile)
	if err != nil {
		return nil, err
	}

	return cm.IndexParams{
		Title:     "Home - " + profile,
		Profile:   profile,
		Reactions: r,
		HashTags:  h,
		Users:     u,
	}.Index(), nil
}

func (l *Logic) ReactionsLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.repo.ReactionNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.repo.HashTagNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) UsersLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.repo.UserNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", name, page),
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) FilesLogic(ctx context.Context, profile string, params Params) (templ.Component, error) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	count, err := l.repo.FileCount(ctx, profile)
	if err != nil {
		return nil, err
	}
	fmt.Println(count)
	p := pg.New(count)
	page := p.Page(params.Page)

	files, err := l.repo.Files(ctx, profile, p)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found")
	}

	hasNext := len(files) >= p.Limit() && p.Next() <= p.Last()
	hasLast := p.Next() < p.Last()

	n := cm.File{
		Title:   fmt.Sprint(page),
		Profile: profile,
		Host:    host,
		Items:   files,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: hasLast,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) NotesLogic(ctx context.Context, profile string, params Params) (templ.Component, error) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	count := 0
	if params.S == "" {
		count, err = l.repo.NoteCount(ctx, profile)
		if err != nil {
			return nil, err
		}
	}

	p := pg.New(count)
	page := p.Page(params.Page)

	notes, err := l.repo.Notes(ctx, profile, params.S, p)
	if err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, fmt.Errorf("note not found")
	}
	title := ""
	if params.S != "" {
		title = fmt.Sprintf("%s - %d", params.S, page)
	} else {
		title = fmt.Sprint(page)
	}

	hasNext := false
	if params.S == "" {
		hasNext = len(notes) >= p.Limit() && p.Next() <= p.Last()
	} else {
		hasNext = len(notes) >= p.Limit()
	}
	hasLast := p.Next() < p.Last()

	n := cm.Note{
		Title:      title,
		Profile:    profile,
		Host:       host,
		SearchPath: fmt.Sprintf("/profiles/%s/notes", profile),
		Items:      notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: hasLast,
		QueryParams: cm.QueryParams{
			Page: params.Page,
			S:    params.S,
		},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) ArchivesLogic(ctx context.Context, profile string) (templ.Component, error) {
	a, err := l.repo.Archives(ctx, profile)
	if err != nil {
		return nil, err
	}
	return cm.ArchiveParams{
		Title:   "Archives",
		Profile: profile,
		Items:   a,
	}.Archive(), nil
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (l *Logic) ArchiveNotesLogic(ctx context.Context, profile, d string, params Params) (templ.Component, error) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	col := ""
	if reym.MatchString(d) {
		col = "strftime('%Y-%m', created_at, 'localtime')"
	} else if reymd.MatchString(d) {
		col = "strftime('%Y-%m-%d', created_at, 'localtime')"
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.repo.ArchiveNotes(ctx, profile, col, d, p)
	if err != nil {
		return nil, err
	}
	title := fmt.Sprintf("%s - %d", d, page)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   title,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) ManageLogic(ctx context.Context) templ.Component {
	p, _ := l.repo.GetProgress()
	var ps []string
	for k := range *l.repo.GetProfiles() {
		ps = append(ps, k)
	}
	// TODO: 進行中の判定これで良いの？
	if p > 0 {
		return cm.ManageStart("Manage")
	}
	return cm.ManageInit("Manage", ps)
}

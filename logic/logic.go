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

func (l *Logic) HomeLogic(ctx context.Context, profile string) (templ.Component, error) {

	if _, ok := l.repo.Config().Profiles[profile]; !ok {
		return nil, fmt.Errorf("invalid profile: %s", profile)
	}

	return cm.IndexParams{
		Title:     "Home",
		Profile:   profile,
		Reactions: l.repo.Reactions(ctx, profile),
		HashTags:  l.repo.HashTags(ctx, profile),
		Users:     l.repo.Users(ctx, profile),
	}.Index(), nil
}

func (l *Logic) ReactionsLogic(ctx context.Context, profile, name string, page int) templ.Component {
	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.ReactionNotes(ctx, profile, name, p)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
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

	return n.WithPages(cp)
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile, name string, page int) templ.Component {
	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.HashTagNotes(ctx, profile, name, p)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
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

	return n.WithPages(cp)
}

func (l *Logic) UsersLogic(ctx context.Context, profile, name string, page int) templ.Component {
	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.UserNotes(ctx, profile, name, p)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", name, page),
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
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

	return n.WithPages(cp)
}

func (l *Logic) FilesLogic(ctx context.Context, profile string, page int) (templ.Component, error) {
	count, err := l.repo.FileCount(ctx, profile)
	if err != nil {
		return nil, err
	}
	fmt.Println(count)
	p := pg.New(count)
	page = p.Page(page)

	files := l.repo.Files(ctx, profile, p)
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found")
	}

	hasNext := len(files) >= p.Limit() && p.Next() <= p.Last()
	hasLast := p.Next() < p.Last()

	n := cm.File{
		Title:   fmt.Sprint(page),
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
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

func (l *Logic) NotesLogic(ctx context.Context, profile string, page int) (templ.Component, error) {
	count, err := l.repo.NoteCount(ctx, profile)
	if err != nil {
		return nil, err
	}
	p := pg.New(count)
	page = p.Page(page)

	notes := l.repo.Notes(ctx, profile, p)
	if len(notes) == 0 {
		return nil, fmt.Errorf("note not found")
	}
	// title := fmt.Sprint(page)

	hasNext := len(notes) >= p.Limit() && p.Next() <= p.Last()
	hasLast := p.Next() < p.Last()

	n := cm.Note{
		Title:   fmt.Sprint(page),
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   notes,
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

func (l *Logic) ArchivesLogic(ctx context.Context, profile string) templ.Component {
	return cm.ArchiveParams{
		Title:   "Archives",
		Profile: profile,
		Items:   l.repo.Archives(ctx, profile),
	}.Archive()
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (l *Logic) ArchiveNotesLogic(ctx context.Context, profile, d string, page int) templ.Component {
	col := ""
	if reym.MatchString(d) {
		col = "strftime('%Y-%m', created_at, 'localtime')"
	} else if reymd.MatchString(d) {
		col = "strftime('%Y-%m-%d', created_at, 'localtime')"
	}

	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.ArchiveNotes(ctx, profile, col, d, p)
	title := fmt.Sprintf("%s - %d", d, page)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   title,
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
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

	return n.WithPages(cp)
}

func (l *Logic) ManageLogic(ctx context.Context) templ.Component {
	p, _ := l.repo.GetProgress()
	var ps []string
	for k := range l.repo.Config().Profiles {
		ps = append(ps, k)
	}
	// TODO: 進行中の判定これで良いの？
	if p > 0 {
		return cm.ManageStart("Manage")
	}
	return cm.ManageInit("Manage", ps)
}

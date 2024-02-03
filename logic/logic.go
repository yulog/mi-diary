package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/infra"
	"github.com/yulog/mi-diary/mi"
	"github.com/yulog/mi-diary/util/pg"
)

type Logic struct {
	repo *infra.Infra
}

func New(r *infra.Infra) *Logic {
	return &Logic{repo: r}
}

func (l *Logic) ProfileLogic(ctx context.Context) templ.Component {
	var ps []string
	for k := range l.repo.Config().Profiles {
		ps = append(ps, k)
	}

	return cm.SelectProfile("Select profile...", ps)
}

func (l *Logic) HomeLogic(ctx context.Context, profile string) templ.Component {

	// TODO: エラーを返すようにする
	if _, ok := l.repo.Config().Profiles[profile]; !ok {
		return nil
	}

	reactions := l.repo.Reactions(ctx, profile)
	tags := l.repo.HashTags(ctx, profile)
	users := l.repo.Users(ctx, profile)

	return cm.Index("Home", profile, cm.Reaction(profile, reactions), cm.HashTag(profile, tags), cm.User(profile, users))
}

func (l *Logic) ReactionsLogic(ctx context.Context, profile, name string) templ.Component {
	notes := l.repo.ReactionNotes(ctx, profile, name)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}

	return n.WithPage()
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile, name string) templ.Component {
	notes := l.repo.HashTagNotes(ctx, profile, name)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}

	return n.WithPage()
}

func (l *Logic) UsersLogic(ctx context.Context, profile, name string) templ.Component {
	notes := l.repo.UserNotes(ctx, profile, name)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}

	return n.WithPage()
}

func (l *Logic) NotesLogic(ctx context.Context, profile string, page int) (templ.Component, error) {
	count, err := l.repo.NoteCount(ctx, profile)
	if err != nil {
		return nil, err
	}
	p := pg.New(count)
	page = p.Page(page)

	notes := l.repo.Notes(ctx, profile, p)
	title := fmt.Sprint(page)

	hasNext := len(notes) >= p.Limit() && p.Next() <= p.Last()
	hasLast := p.Next() < p.Last()

	n := cm.Note{
		Title:   title,
		Profile: profile,
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
	archives := l.repo.Archives(ctx, profile)

	return cm.Archive("Archives", profile, archives)
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

func (l *Logic) SettingsLogic(ctx context.Context, profile string) templ.Component {

	return cm.Settings("settings", profile)
}

func (l *Logic) SettingsReactionsLogic(ctx context.Context, profile, id string) {
	body := map[string]any{
		"i":      l.repo.Config().Profiles[profile].I,
		"limit":  20,
		"userId": l.repo.Config().Profiles[profile].UserId,
	}
	if id != "" {
		body["untilId"] = id
	}
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/users/reactions", l.repo.Config().Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	l.repo.Insert(ctx, profile, resp)
}

func (l *Logic) SettingsEmojisLogic(ctx context.Context, profile, name string) {
	body := map[string]any{
		"name": name,
	}
	// if id != "" {
	// 	body["untilId"] = id
	// }
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/emoji", l.repo.Config().Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	l.repo.InsertEmoji(ctx, profile, resp)
}

package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/a-h/templ"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/mi"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/util/pg"
)

type Logic struct {
	app *app.App
}

func New(a *app.App) *Logic {
	return &Logic{app: a}
}

func (l Logic) ProfileLogic(ctx context.Context) templ.Component {
	var ps []string
	for k := range l.app.Config.Profiles {
		ps = append(ps, k)
	}

	return cm.SelectProfile("Select profile...", ps)
}

func (l Logic) HomeLogic(ctx context.Context, profile string) templ.Component {
	var reactions []model.Reaction
	l.app.DB(profile).
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(ctx)
	var tags []model.HashTag
	l.app.DB(profile).
		NewSelect().
		Model(&tags).
		Order("count DESC").
		Scan(ctx)
	var users []model.User
	l.app.DB(profile).
		NewSelect().
		Model(&users).
		Order("count DESC").
		Scan(ctx)

	return cm.Index("Home", profile, cm.Reaction(profile, reactions), cm.HashTag(profile, tags), cm.User(profile, users))
}

func (l Logic) ReactionsLogic(ctx context.Context, profile, name string) templ.Component {
	var notes []model.Note
	l.app.DB(profile).
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Order("created_at DESC").
		Scan(ctx)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}

	return n.WithPage()
}

func (l Logic) HashTagsLogic(ctx context.Context, profile, name string) templ.Component {
	var notes []model.Note
	l.app.DB(profile).
		NewSelect().
		Model((*model.NoteToTag)(nil)).
		// 必要な列だけ選択して、不要な列をなくす
		Relation("Note", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "user_id", "reaction_name", "text")
		}).
		Relation("HashTag", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		Column("").
		Where("hash_tag.text = ?", name).
		Order("created_at DESC").
		Scan(ctx, &notes)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}

	return n.WithPage()
}

func (l Logic) UsersLogic(ctx context.Context, profile, name string) templ.Component {
	var notes []model.Note
	l.app.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Where("user.name = ?", name).
		Order("created_at DESC").
		Scan(ctx)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}

	return n.WithPage()
}

func (l Logic) NotesLogic(ctx context.Context, profile string, page int) (templ.Component, error) {
	count, err := l.app.DB(profile).
		NewSelect().
		Model((*model.Note)(nil)).
		Count(ctx)
	if err != nil {
		return nil, err
	}
	p := pg.New(count)
	page = p.Page(page)

	var notes []model.Note
	l.app.DB(profile).
		NewSelect().
		Model(&notes).
		// Relation("User").
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
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

func (l Logic) ArchivesLogic(ctx context.Context, profile string) templ.Component {
	var archives []model.Archive
	l.app.DB(profile).
		NewSelect().
		Model((*model.Day)(nil)).
		Relation("Month", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		ColumnExpr("d.ym as ym, month.count as ym_count, d.ymd as ymd, d.count as ymd_count").
		Order("ym DESC", "ymd DESC").
		Scan(ctx, &archives)

	return cm.Archive("Archives", profile, archives)
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (l Logic) ArchiveNotesLogic(ctx context.Context, profile, d string, page int) templ.Component {
	col := ""
	if reym.MatchString(d) {
		col = "strftime('%Y-%m', created_at, 'localtime')"
	} else if reymd.MatchString(d) {
		col = "strftime('%Y-%m-%d', created_at, 'localtime')"
	}

	p := pg.New(0)
	page = p.Page(page)

	var notes []model.Note
	l.app.DB(profile).
		NewSelect().
		Model(&notes).
		Where(col+" = ?", d). // 条件指定に関数適用した列を使う
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
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

func (l Logic) GetReactions(ctx context.Context, profile, id string) {
	body := map[string]any{
		"i":      l.app.Config.Profiles[profile].I,
		"limit":  20,
		"userId": l.app.Config.Profiles[profile].UserId,
	}
	if id != "" {
		body["untilId"] = id
	}
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/users/reactions", l.app.Config.Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	l.app.Insert(ctx, profile, resp)
}

func (l Logic) GetEmojiOne(ctx context.Context, profile, name string) {
	body := map[string]any{
		"name": name,
	}
	// if id != "" {
	// 	body["untilId"] = id
	// }
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/emoji", l.app.Config.Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	l.app.InsertEmoji(ctx, profile, resp)
}

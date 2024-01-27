package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/mi"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/server/pg"
)

// ProfileHandler は / のハンドラ
func (srv *Server) ProfileHandler(c echo.Context) error {
	var ps []string
	for k := range srv.app.Config.Profiles {
		ps = append(ps, k)
	}
	return renderer(c, cm.SelectProfile("Select profile...", ps))
}

// HomeHandler は /:profile のハンドラ
func (srv *Server) HomeHandler(c echo.Context) error {
	profile := c.Param("profile")
	var reactions []model.Reaction
	srv.app.DB().
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(c.Request().Context())
	var tags []model.HashTag
	srv.app.DB().
		NewSelect().
		Model(&tags).
		Order("count DESC").
		Scan(c.Request().Context())
	var users []model.User
	srv.app.DB().
		NewSelect().
		Model(&users).
		Order("count DESC").
		Scan(c.Request().Context())
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, cm.Index("Home", profile, cm.Reaction(profile, reactions), cm.HashTag(profile, tags), cm.User(profile, users)))
}

// ReactionsHandler は /reactions/:name のハンドラ
func (srv *Server) ReactionsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Order("created_at DESC").
		Scan(c.Request().Context())
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}
	return renderer(c, n.WithPage())
}

// HashTagsHandler は /hashtags/:name のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
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
		Scan(c.Request().Context(), &notes)
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}
	return renderer(c, n.WithPage())
}

// UsersHandler は /users/:name のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Relation("User").
		Where("user.name = ?", name).
		Order("created_at DESC").
		Scan(c.Request().Context())
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}
	return renderer(c, n.WithPage())
}

// NotesHandler は /notes のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	count, err := srv.app.DB().
		NewSelect().
		Model((*model.Note)(nil)).
		Count(c.Request().Context())
	if err != nil {
		return err
	}
	p := pg.New(&c, count)
	page, err := p.Page()
	if err != nil {
		return err
	}

	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		// Relation("User").
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(c.Request().Context())
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
	return renderer(c, n.WithPages(cp))
}

// ArchivesHandler は /archives のハンドラ
func (srv *Server) ArchivesHandler(c echo.Context) error {
	profile := c.Param("profile")
	var archives []model.Archive
	srv.app.DB().
		NewSelect().
		Model((*model.Day)(nil)).
		Relation("Month", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		ColumnExpr("d.ym as ym, month.count as ym_count, d.ymd as ymd, d.count as ymd_count").
		Order("ym DESC", "ymd DESC").
		Scan(c.Request().Context(), &archives)
	return renderer(c, cm.Archive("Archives", profile, archives))
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// ArchiveNotesHandler は /archives/:date のハンドラ
func (srv *Server) ArchiveNotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	d := c.Param("date")
	col := ""
	if reym.MatchString(d) {
		col = "strftime('%Y-%m', created_at, 'localtime')"
	} else if reymd.MatchString(d) {
		col = "strftime('%Y-%m-%d', created_at, 'localtime')"
	}

	p := pg.New(&c, 0)
	page, err := p.Page()
	if err != nil {
		return err
	}

	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Where(col+" = ?", d). // 条件指定に関数適用した列を使う
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(c.Request().Context())
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
	return renderer(c, n.WithPages(cp))
}

// SettingsHandler は /settings のハンドラ
func (srv *Server) SettingsHandler(c echo.Context) error {
	profile := c.Param("profile")
	return renderer(c, cm.Settings("settings", profile))
}

// SettingsReactionsHandler は /settings/reactions のハンドラ
func (srv *Server) SettingsReactionsHandler(c echo.Context) error {
	profile := c.Param("profile")
	id := c.FormValue("note-id")
	body := map[string]any{
		"i":      srv.app.Config.Profiles[profile].I,
		"limit":  20,
		"userId": srv.app.Config.Profiles[profile].UserId,
	}
	if id != "" {
		body["untilId"] = id
	}
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/users/reactions", srv.app.Config.Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	app.Insert(c.Request().Context(), resp)
	return c.HTML(http.StatusOK, id)
}

// SettingsEmojisHandler は /settings/emojis のハンドラ
func (srv *Server) SettingsEmojisHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.FormValue("emoji-name")
	body := map[string]any{
		"name": name,
	}
	// if id != "" {
	// 	body["untilId"] = id
	// }
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/emoji", srv.app.Config.Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	app.InsertEmoji(c.Request().Context(), resp)
	return c.HTML(http.StatusOK, name)
}

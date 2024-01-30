package server

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	cm "github.com/yulog/mi-diary/components"
)

// ProfileHandler は / のハンドラ
func (srv *Server) ProfileHandler(c echo.Context) error {
	/*var ps []string
	for k := range srv.app.Config.Profiles {
		ps = append(ps, k)
	}*/
	return renderer(c, srv.logic.ProfileLogic(c.Request().Context()))
}

// HomeHandler は /:profile のハンドラ
func (srv *Server) HomeHandler(c echo.Context) error {
	profile := c.Param("profile")
	/*var reactions []model.Reaction
	srv.app.DB(profile).
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(c.Request().Context())
	var tags []model.HashTag
	srv.app.DB(profile).
		NewSelect().
		Model(&tags).
		Order("count DESC").
		Scan(c.Request().Context())
	var users []model.User
	srv.app.DB(profile).
		NewSelect().
		Model(&users).
		Order("count DESC").
		Scan(c.Request().Context())*/
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, srv.logic.HomeLogic(c.Request().Context(), profile))
}

// ReactionsHandler は /reactions/:name のハンドラ
func (srv *Server) ReactionsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	/*var notes []model.Note
	srv.app.DB(profile).
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Order("created_at DESC").
		Scan(c.Request().Context())
	n := cm.Note{
		Title:   name,
		Profile: profile,
		Items:   notes,
	}*/
	return renderer(c, srv.logic.ReactionsLogic(c.Request().Context(), profile, name))
}

// HashTagsHandler は /hashtags/:name のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	/*var notes []model.Note
	srv.app.DB(profile).
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
	}*/
	return renderer(c, srv.logic.HashTagsLogic(c.Request().Context(), profile, name))
}

// UsersHandler は /users/:name のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	/*var notes []model.Note
	srv.app.DB(profile).
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
	}*/
	return renderer(c, srv.logic.UsersLogic(c.Request().Context(), profile, name))
}

// NotesHandler は /notes のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	var p int
	if err := echo.QueryParamsBinder(c).
		Int("page", &p).
		BindError(); err != nil {
		return err
	}
	/*count, err := srv.app.DB(profile).
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
	srv.app.DB(profile).
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
	}*/
	com, err := srv.logic.NotesLogic(c.Request().Context(), profile, p)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ArchivesHandler は /archives のハンドラ
func (srv *Server) ArchivesHandler(c echo.Context) error {
	profile := c.Param("profile")
	/*var archives []model.Archive
	srv.app.DB(profile).
		NewSelect().
		Model((*model.Day)(nil)).
		Relation("Month", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		ColumnExpr("d.ym as ym, month.count as ym_count, d.ymd as ymd, d.count as ymd_count").
		Order("ym DESC", "ymd DESC").
		Scan(c.Request().Context(), &archives)*/
	return renderer(c, srv.logic.ArchivesLogic(c.Request().Context(), profile))
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// ArchiveNotesHandler は /archives/:date のハンドラ
func (srv *Server) ArchiveNotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	d := c.Param("date")
	var p int
	if err := echo.QueryParamsBinder(c).
		Int("page", &p).
		BindError(); err != nil {
		return err
	}
	/*col := ""
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
	srv.app.DB(profile).
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
	}*/

	return renderer(c, srv.logic.ArchiveNotesLogic(c.Request().Context(), profile, d, p))
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
	/*body := map[string]any{
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
	srv.app.Insert(c.Request().Context(), profile, resp)*/
	srv.logic.GetReactions(c.Request().Context(), profile, id)
	return c.HTML(http.StatusOK, id)
}

// SettingsEmojisHandler は /settings/emojis のハンドラ
func (srv *Server) SettingsEmojisHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.FormValue("emoji-name")
	/*body := map[string]any{
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
	srv.app.InsertEmoji(c.Request().Context(), profile, resp)*/
	srv.logic.GetEmojiOne(c.Request().Context(), profile, name)
	return c.HTML(http.StatusOK, name)
}

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/mi"
	"github.com/yulog/mi-diary/model"
)

// IndexHandler は / のハンドラ
func (srv *Server) IndexHandler(c echo.Context) error {
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
	return renderer(c, cm.Index("index", cm.Reaction(reactions), cm.HashTag(tags), cm.User(users)))
}

// ReactionsHandler は /reactions/:name のハンドラ
func (srv *Server) ReactionsHandler(c echo.Context) error {
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Scan(c.Request().Context())
	return renderer(c, cm.Note(name, notes))
}

// HashTagsHandler は /hashtags/:name のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
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
		Scan(c.Request().Context(), &notes)
	return renderer(c, cm.Note(name, notes))
}

// UsersHandler は /users/:name のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Relation("User").
		Where("user.name = ?", name).
		Scan(c.Request().Context())
	return renderer(c, cm.Note(name, notes))
}

// NotesHandler は /notes/:page のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	var page = 1
	err := echo.QueryParamsBinder(c).
		Int("page", &page).
		BindError()
	// page, err = strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return err
	}
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := limit * (page - 1)
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		// Relation("User").
		Limit(limit).
		Offset(offset).
		Scan(c.Request().Context())
	title := fmt.Sprint(page)
	prev := page - 1
	next := page + 1
	if len(notes) < limit {
		next = 0
	}
	return renderer(c, cm.NoteWithPages(title, notes, page, prev, next))
}

// SettingsHandler は /settings のハンドラ
func (srv *Server) SettingsHandler(c echo.Context) error {
	return renderer(c, cm.Settings("settings"))
}

// SettingsReactionsHandler は /settings/reactions のハンドラ
func (srv *Server) SettingsReactionsHandler(c echo.Context) error {
	id := c.FormValue("note-id")
	body := map[string]any{
		"i":      srv.app.Config.I,
		"limit":  20,
		"userId": srv.app.Config.UserId,
	}
	if id != "" {
		body["untilId"] = id
	}
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/users/reactions", srv.app.Config.Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	app.Insert(c.Request().Context(), resp)
	return c.HTML(http.StatusOK, id)
}

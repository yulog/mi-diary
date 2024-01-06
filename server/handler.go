package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/model"
)

// IndexHandler は/のハンドラ
func (srv *Server) IndexHandler(c echo.Context) error {
	var reactions []model.Reaction
	srv.app.DB().
		NewSelect().
		Model(&reactions).
		Scan(c.Request().Context())
	var tags []model.HashTag
	srv.app.DB().
		NewSelect().
		Model(&tags).
		Scan(c.Request().Context())
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, cm.Index("index", cm.Reaction(reactions), cm.HashTag(tags)))
}

// ReactionsHandler は/reactions/:nameのハンドラ
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

// HashTagsHandler は/hashtags/:nameのハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model((*model.NoteToTag)(nil)).
		// 必要な列だけ選択して、不要な列をなくす
		Relation("Note", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "user_id", "reaction_name")
		}).
		Relation("HashTag", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		Column("").
		Where("hash_tag.text = ?", name).
		Scan(c.Request().Context(), &notes)
	return renderer(c, cm.Note(name, notes))
}

// SettingsHandler は/settingsのハンドラ
func (srv *Server) SettingsHandler(c echo.Context) error {
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, cm.Settings("settings"))
}

// SettingsReactionsHandler は/settings/reactionsのハンドラ
func (srv *Server) SettingsReactionsHandler(c echo.Context) error {
	app.Insert(c.Request().Context())
	return c.HTML(http.StatusOK, "OK")
	// return renderer(c, cm.Settings("settings"))
}

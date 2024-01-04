package server

import (
	"github.com/a-h/templ"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/model"

	"github.com/labstack/echo"
)

type Server struct {
	app *app.App
}

func New(a *app.App) *Server {
	return &Server{app: a}
}

func renderer(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (srv *Server) GetIndex(c echo.Context) error {
	var reactions []model.Reaction
	srv.app.DB().
		NewSelect().
		Model(&reactions).
		Scan(c.Request().Context())
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, cm.Index("index", cm.Reaction(reactions)))
}

func (srv *Server) GetReactions(c echo.Context) error {
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Scan(c.Request().Context())
	return renderer(c, cm.Note(notes))
}

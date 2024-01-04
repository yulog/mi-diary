package main

import (
	"github.com/a-h/templ"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func (srv *Server) getIndex(c echo.Context) error {
	var reactions []model.Reaction
	srv.app.DB().
		NewSelect().
		Model(&reactions).
		Scan(c.Request().Context())
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, cm.Index("index", cm.Reaction(reactions)))
}

func (srv *Server) getReactions(c echo.Context) error {
	name := c.Param("name")
	var notes []model.Note
	srv.app.DB().
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Scan(c.Request().Context())
	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, cm.Note(notes))
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	app := app.New()
	srv := New(app)
	// e.Validator = &Validator{validator: validator.New()}
	e.GET("/", srv.getIndex)
	e.GET("/reactions/:name", srv.getReactions)
	e.Logger.Fatal(e.Start(":1323"))
}

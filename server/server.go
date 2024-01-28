package server

import (
	"github.com/a-h/templ"
	"github.com/yulog/mi-diary/app"

	"github.com/labstack/echo/v4"
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

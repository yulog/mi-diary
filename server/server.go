package server

import (
	"github.com/a-h/templ"
	"github.com/yulog/mi-diary/logic"

	"github.com/labstack/echo/v4"
)

type Server struct {
	logic *logic.Logic
}

func New(l *logic.Logic) *Server {
	return &Server{logic: l}
}

func renderer(c echo.Context, cmp templ.Component) error {
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

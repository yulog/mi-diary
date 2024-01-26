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

// type Pager struct {
// 	Page int
// 	Prev int
// 	Next int
// 	Last int

// 	Limit  int
// 	Offset int
// }

// func pager(c echo.Context, part, total int) (Pager, error) {
// 	var page = 1
// 	err := echo.QueryParamsBinder(c).
// 		Int("page", &page).
// 		BindError()
// 	if err != nil {
// 		return Pager{}, err
// 	}
// 	if page < 1 {
// 		page = 1
// 	}
// 	limit := 10
// 	offset := limit * (page - 1)

// 	last := int(math.Ceil(float64(total) / float64(limit)))

// 	prev := page - 1
// 	next := page + 1
// 	if part < limit || next > last {
// 		next = 0
// 	}
// 	if next == last {
// 		last = 0
// 	}
// 	return Pager{
// 		Page: page,
// 		Prev: prev,
// 		Next: next,
// 		Last: last,

// 		Limit:  limit,
// 		Offset: offset,
// 	}, nil
// }

// func page(c echo.Context) (int, error) {
// 	var page = 1
// 	if err := echo.QueryParamsBinder(c).
// 		Int("page", &page).
// 		BindError(); err != nil {
// 		return page, err
// 	}
// 	if page < 1 {
// 		page = 1
// 	}
// 	return page, nil
// }

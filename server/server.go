package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/yulog/mi-diary/logic"

	"github.com/labstack/echo/v4"
)

type Server struct {
	logic *logic.Logic
}

func New(l *logic.Logic) *Server {
	return &Server{logic: l}
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func renderer(c echo.Context, cmp templ.Component) error {
	// https://github.com/a-h/templ/blob/067cc686cd1e44cd0d3b6669a52e24ef115ccc5a/examples/integration-echo/main.go#L17
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := cmp.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, buf.String())
}

func page(c echo.Context, p *int) error {
	if err := echo.QueryParamsBinder(c).
		Int("page", p).
		BindError(); err != nil {
		return err
	}
	return nil
}

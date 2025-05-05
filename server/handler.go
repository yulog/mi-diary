package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yulog/mi-diary/presenter"
)

type Params struct {
	Profile string `param:"profile"`
	Name    string `param:"name"`
	Date    string `param:"date"`
	Page    int    `query:"page"`
	S       string `query:"s"`
	Color   string `query:"color"`
	Partial bool   `query:"partial"`
	SortBy  string `query:"sort"`
}

type Callback struct {
	Host      string `param:"host"`
	SessionID string `query:"session"`
}

type Job struct {
	Profile string `form:"profile" validate:"required"`
	Type    int    `form:"job-type" validate:"required"`
	ID      string `form:"id"`
}

type Profiles struct {
	ServerURL string `form:"server-url"`
}

// RootHandler は / のハンドラ
func (srv *Server) RootHandler(c echo.Context) error {
	return renderer(c, presenter.SelectProfilePresentation(c, srv.logic.SelectProfileLogic(c.Request().Context())))
}

// NewProfilesHandler は /profiles のハンドラ
func (srv *Server) NewProfilesHandler(c echo.Context) error {

	return renderer(c, presenter.AddProfilePresentation(c, srv.logic.NewProfileLogic(c.Request().Context())))
}

// AddProfileHandler は /profiles のハンドラ
func (srv *Server) AddProfileHandler(c echo.Context) error {
	var params Profiles
	if err := c.Bind(&params); err != nil {
		return err
	}

	authURL, err := srv.logic.AddProfileLogic(c.Request().Context(), params.ServerURL)
	if err != nil {
		return err
	}

	c.Response().Header().Set("hx-redirect", authURL)

	return c.NoContent(http.StatusOK)
}

// CallbackHandler は /callback/:host のハンドラ
func (srv *Server) CallbackHandler(c echo.Context) error {
	var callback Callback
	if err := c.Bind(&callback); err != nil {
		return err
	}

	err := srv.logic.CallbackLogic(c.Request().Context(), callback.Host, callback.SessionID)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

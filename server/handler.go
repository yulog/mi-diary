//go:generate go tool govalid .

package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	cm "github.com/yulog/mi-diary/components"
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
	//govalid:required
	Profile string `form:"profile"`
	//govalid:required
	Type int    `form:"job-type"`
	ID   string `form:"id"`
}

type Profiles struct {
	ServerURL string `form:"server-url"`
}

// RootHandler は / のハンドラ
func (srv *Server) RootHandler(c *echo.Context) error {
	return renderer(c, cm.SelectProfile("Select profile...", srv.logic.ConfigRepo.GetProfilesSortedKey()))
}

// NewProfilesHandler は /profiles のハンドラ
func (srv *Server) NewProfilesHandler(c *echo.Context) error {
	return renderer(c, cm.AddProfile("New Profile"))
}

// AddProfileHandler は /profiles のハンドラ
func (srv *Server) AddProfileHandler(c *echo.Context) error {
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
func (srv *Server) CallbackHandler(c *echo.Context) error {
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

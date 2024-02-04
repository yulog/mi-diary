package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RootHandler は / のハンドラ
func (srv *Server) RootHandler(c echo.Context) error {
	return renderer(c, srv.logic.SelectProfileLogic(c.Request().Context()))
}

// HomeHandler は /:profile のハンドラ
func (srv *Server) HomeHandler(c echo.Context) error {
	profile := c.Param("profile")

	// TODO: logic の返り値をチェックして
	// echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// logic はDTOを返すようにする？
	// エラーでなければ、DTOのメソッドでcomponentを作る

	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	return renderer(c, srv.logic.HomeLogic(c.Request().Context(), profile))
}

// ReactionsHandler は /reactions/:name のハンドラ
func (srv *Server) ReactionsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")

	return renderer(c, srv.logic.ReactionsLogic(c.Request().Context(), profile, name))
}

// HashTagsHandler は /hashtags/:name のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")

	return renderer(c, srv.logic.HashTagsLogic(c.Request().Context(), profile, name))
}

// UsersHandler は /users/:name のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")

	return renderer(c, srv.logic.UsersLogic(c.Request().Context(), profile, name))
}

// NotesHandler は /notes のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	var p int
	if err := echo.QueryParamsBinder(c).
		Int("page", &p).
		BindError(); err != nil {
		return err
	}

	com, err := srv.logic.NotesLogic(c.Request().Context(), profile, p)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ArchivesHandler は /archives のハンドラ
func (srv *Server) ArchivesHandler(c echo.Context) error {
	profile := c.Param("profile")

	return renderer(c, srv.logic.ArchivesLogic(c.Request().Context(), profile))
}

// ArchiveNotesHandler は /archives/:date のハンドラ
func (srv *Server) ArchiveNotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	d := c.Param("date")
	var p int
	if err := echo.QueryParamsBinder(c).
		Int("page", &p).
		BindError(); err != nil {
		return err
	}

	return renderer(c, srv.logic.ArchiveNotesLogic(c.Request().Context(), profile, d, p))
}

// SettingsHandler は /settings のハンドラ
func (srv *Server) SettingsHandler(c echo.Context) error {
	profile := c.Param("profile")

	return renderer(c, srv.logic.SettingsLogic(c.Request().Context(), profile))
}

// SettingsReactionsHandler は /settings/reactions のハンドラ
func (srv *Server) SettingsReactionsHandler(c echo.Context) error {
	profile := c.Param("profile")
	id := c.FormValue("note-id")

	srv.logic.SettingsReactionsLogic(c.Request().Context(), profile, id)
	return c.HTML(http.StatusOK, id)
}

// SettingsEmojisHandler は /settings/emojis のハンドラ
func (srv *Server) SettingsEmojisHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.FormValue("emoji-name")

	srv.logic.SettingsEmojisLogic(c.Request().Context(), profile, name)
	return c.HTML(http.StatusOK, name)
}

// NewProfilesHandler は /profiles のハンドラ
func (srv *Server) NewProfilesHandler(c echo.Context) error {

	return renderer(c, srv.logic.NewProfileLogic(c.Request().Context()))
}

// AddProfileHandler は /profiles のハンドラ
func (srv *Server) AddProfileHandler(c echo.Context) error {
	server := c.FormValue("server-url")

	authURL := srv.logic.AddProfileLogic(c.Request().Context(), server)

	c.Response().Header().Set("hx-redirect", authURL)

	return c.NoContent(http.StatusOK)
}

// CallbackHandler は /callback/:host のハンドラ
func (srv *Server) CallbackHandler(c echo.Context) error {
	host := c.Param("host")
	s := c.QueryParam("session")

	srv.logic.CallbackLogic(c.Request().Context(), host, s)

	return c.Redirect(http.StatusFound, "/")
}

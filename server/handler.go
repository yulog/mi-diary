package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yulog/mi-diary/app"
)

type Job struct {
	Profile string `form:"profile" validate:"required"`
	Type    int    `form:"job-type" validate:"required"`
	ID      string `form:"id"`
}

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
	com, err := srv.logic.HomeLogic(c.Request().Context(), profile)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ReactionsHandler は /reactions/:name のハンドラ
func (srv *Server) ReactionsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	var p int
	if err := page(c, &p); err != nil {
		return err
	}

	com, err := srv.logic.ReactionsLogic(c.Request().Context(), profile, name, p)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// HashTagsHandler は /hashtags/:name のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	var p int
	if err := page(c, &p); err != nil {
		return err
	}

	com, err := srv.logic.HashTagsLogic(c.Request().Context(), profile, name, p)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// UsersHandler は /users/:name のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	profile := c.Param("profile")
	name := c.Param("name")
	var p int
	if err := page(c, &p); err != nil {
		return err
	}

	com, err := srv.logic.UsersLogic(c.Request().Context(), profile, name, p)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// FilesHandler は /files のハンドラ
func (srv *Server) FilesHandler(c echo.Context) error {
	profile := c.Param("profile")
	var p int
	if err := page(c, &p); err != nil {
		return err
	}

	com, err := srv.logic.FilesLogic(c.Request().Context(), profile, p)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// NotesHandler は /notes のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	var p int
	if err := page(c, &p); err != nil {
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

	com, err := srv.logic.ArchivesLogic(c.Request().Context(), profile)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ArchiveNotesHandler は /archives/:date のハンドラ
func (srv *Server) ArchiveNotesHandler(c echo.Context) error {
	profile := c.Param("profile")
	d := c.Param("date")
	var p int
	if err := page(c, &p); err != nil {
		return err
	}

	com, err := srv.logic.ArchiveNotesLogic(c.Request().Context(), profile, d, p)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// ManageHandler は /manage のハンドラ
func (srv *Server) ManageHandler(c echo.Context) error {

	return renderer(c, srv.logic.ManageLogic(c.Request().Context()))
}

// JobStartHandler は /job/start のハンドラ
func (srv *Server) JobStartHandler(c echo.Context) error {
	j := new(Job)
	if err := c.Bind(j); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(j); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	job := app.Job{
		Profile: j.Profile,
		Type:    app.JobType(j.Type),
		ID:      j.ID,
	}

	return renderer(c, srv.logic.JobStartLogic(c.Request().Context(), job))
}

// JobProgressHandler は /job/progress のハンドラ
func (srv *Server) JobProgressHandler(c echo.Context) error {
	_, d, com := srv.logic.JobProgressLogic(c.Request().Context())

	if d {
		c.Response().Header().Set("hx-trigger", "done")
	}

	return renderer(c, com)
}

// JobHandler は /job のハンドラ
func (srv *Server) JobHandler(c echo.Context) error {
	profile := c.Param("profile")
	com := srv.logic.JobLogic(c.Request().Context(), profile)

	return renderer(c, com)
}

// NewProfilesHandler は /profiles のハンドラ
func (srv *Server) NewProfilesHandler(c echo.Context) error {

	return renderer(c, srv.logic.NewProfileLogic(c.Request().Context()))
}

// AddProfileHandler は /profiles のハンドラ
func (srv *Server) AddProfileHandler(c echo.Context) error {
	server := c.FormValue("server-url")

	authURL, err := srv.logic.AddProfileLogic(c.Request().Context(), server)
	if err != nil {
		return err
	}

	c.Response().Header().Set("hx-redirect", authURL)

	return c.NoContent(http.StatusOK)
}

// CallbackHandler は /callback/:host のハンドラ
func (srv *Server) CallbackHandler(c echo.Context) error {
	host := c.Param("host")
	s := c.QueryParam("session")

	err := srv.logic.CallbackLogic(c.Request().Context(), host, s)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

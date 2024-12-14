package server

import (
	"github.com/labstack/echo/v4"
	"github.com/yulog/mi-diary/logic"
)

// HomeHandler は /:profile のハンドラ
func (srv *Server) HomeHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	// TODO: logic の返り値をチェックして
	// echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// logic はDTOを返すようにする？
	// エラーでなければ、DTOのメソッドでcomponentを作る

	// return c.HTML(http.StatusOK, fmt.Sprint(reactions))
	com, err := srv.logic.HomeLogic(c.Request().Context(), params.Profile)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ReactionHandler は /:profile/reactions/:name のハンドラ
func (srv *Server) ReactionHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := logic.Params{
		Page: params.Page,
	}

	com, err := srv.logic.ReactionNotesLogic(c.Request().Context(), params.Profile, params.Name, params2)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// HashTagsHandler は /:profile/hashtags のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	com, err := srv.logic.HashTagsLogic(c.Request().Context(), params.Profile)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// HashTagHandler は /:profile/hashtags/:name のハンドラ
func (srv *Server) HashTagHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := logic.Params{
		Page: params.Page,
	}

	com, err := srv.logic.HashTagNotesLogic(c.Request().Context(), params.Profile, params.Name, params2)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// UsersHandler は /:profile/users のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	com, err := srv.logic.UsersLogic(c.Request().Context(), params.Profile, params.Partial, params.SortBy)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// UserHandler は /:profile/users/:name のハンドラ
func (srv *Server) UserHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := logic.Params{
		Page: params.Page,
	}

	com, err := srv.logic.UserLogic(c.Request().Context(), params.Profile, params.Name, params2)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

// FilesHandler は /:profile/files のハンドラ
func (srv *Server) FilesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := logic.Params{
		Page:  params.Page,
		Color: params.Color,
	}

	com, err := srv.logic.FilesLogic(c.Request().Context(), params.Profile, params2)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// NotesHandler は /:profile/notes のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := logic.Params{
		Page: params.Page,
		S:    params.S,
	}

	com, err := srv.logic.NotesLogic(c.Request().Context(), params.Profile, params2)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ArchivesHandler は /:profile/archives のハンドラ
func (srv *Server) ArchivesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	com, err := srv.logic.ArchivesLogic(c.Request().Context(), params.Profile)
	if err != nil {
		return err
	}
	return renderer(c, com)
}

// ArchiveNotesHandler は /:profile/archives/:date のハンドラ
func (srv *Server) ArchiveNotesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := logic.Params{
		Page: params.Page,
	}

	com, err := srv.logic.ArchiveNotesLogic(c.Request().Context(), params.Profile, params.Date, params2)
	if err != nil {
		return err
	}

	return renderer(c, com)
}

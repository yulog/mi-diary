package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/yulog/mi-diary/internal/shared"
	"github.com/yulog/mi-diary/presenter"
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
	out, err := srv.logic.HomeLogic(c.Request().Context(), params.Profile)
	if err != nil {
		return err
	}
	return renderer(c, presenter.IndexPresentation(c, out))
}

// ReactionHandler は /:profile/reactions/:name のハンドラ
func (srv *Server) ReactionHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := shared.QueryParams{
		Page: params.Page,
	}

	out, err := srv.logic.ReactionNotesLogic(c.Request().Context(), params.Profile, params.Name, params2)
	if err != nil {
		return err
	}

	return renderer(c, presenter.NoteWithPagesPresentation(c, out))
}

// HashTagsHandler は /:profile/hashtags のハンドラ
func (srv *Server) HashTagsHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	out, err := srv.logic.HashTagsLogic(c.Request().Context(), params.Profile)
	if err != nil {
		return err
	}

	return renderer(c, presenter.HashTagPresentation(c, out))
}

// HashTagHandler は /:profile/hashtags/:name のハンドラ
func (srv *Server) HashTagHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := shared.QueryParams{
		Page: params.Page,
	}

	out, err := srv.logic.HashTagNotesLogic(c.Request().Context(), params.Profile, params.Name, params2)
	if err != nil {
		return err
	}

	return renderer(c, presenter.NoteWithPagesPresentation(c, out))
}

// UsersHandler は /:profile/users のハンドラ
func (srv *Server) UsersHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	out, err := srv.logic.UsersLogic(c.Request().Context(), params.Profile, params.Partial, params.SortBy)
	if err != nil {
		return err
	}

	return renderer(c, presenter.UserPresentation(c, out))
}

// UserHandler は /:profile/users/:name のハンドラ
func (srv *Server) UserHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := shared.QueryParams{
		Page: params.Page,
	}

	out, err := srv.logic.UserLogic(c.Request().Context(), params.Profile, params.Name, params2)
	if err != nil {
		return err
	}

	return renderer(c, presenter.NoteWithPagesPresentation(c, out))
}

// FilesHandler は /:profile/files のハンドラ
func (srv *Server) FilesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := shared.QueryParams{
		Page:  params.Page,
		Color: params.Color,
	}

	out, err := srv.logic.FilesLogic(c.Request().Context(), params.Profile, params2)
	if err != nil {
		return err
	}
	return renderer(c, presenter.FileWithPagesPresentation(c, out))
}

// NotesHandler は /:profile/notes のハンドラ
func (srv *Server) NotesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := shared.QueryParams{
		Page: params.Page,
		S:    params.S,
	}

	out, err := srv.logic.NotesLogic(c.Request().Context(), params.Profile, params2)
	if err != nil {
		return err
	}
	return renderer(c, presenter.NoteWithPagesPresentation(c, out))
}

// ArchivesHandler は /:profile/archives のハンドラ
func (srv *Server) ArchivesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	out, err := srv.logic.ArchivesLogic(c.Request().Context(), params.Profile)
	if err != nil {
		return err
	}
	return renderer(c, presenter.ArchivesPresentation(c, out))
}

// ArchiveNotesHandler は /:profile/archives/:date のハンドラ
func (srv *Server) ArchiveNotesHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}
	params2 := shared.QueryParams{
		Page: params.Page,
	}

	out, err := srv.logic.ArchiveNotesLogic(c.Request().Context(), params.Profile, params.Date, params2)
	if err != nil {
		return err
	}

	return renderer(c, presenter.NoteWithPagesPresentation(c, out))
}

// EmojiHandler は /:profile/emojis/:name のハンドラ
func (srv *Server) EmojiHandler(c echo.Context) error {
	var params Params
	if err := c.Bind(&params); err != nil {
		return err
	}

	resp, err := srv.logic.CacheLogic(c.Request().Context(), params.Profile, params.Name)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Response.Close()
	defer resp.DoCache()

	c.Stream(resp.StatusCode, resp.ContentType, resp.Response)
	return nil
}

package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (srv *Server) NewRouter() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Validator = NewValidator()

	e.GET("/", srv.RootHandler)
	e.GET("/callback/:host", srv.CallbackHandler)
	e.GET("/manage", srv.ManageHandler)

	job := e.Group("/job")
	job.GET("", srv.JobHandler)
	job.GET("/progress", srv.JobProgressHandler)
	job.POST("/start", srv.JobStartHandler)

	profiles := e.Group("/profiles")
	profiles.GET("", srv.NewProfilesHandler)
	profiles.POST("", srv.AddProfileHandler)

	profile := profiles.Group("/:profile")
	profile.GET("", srv.HomeHandler)
	profile.GET("/reactions/:name", srv.ReactionHandler)
	profile.GET("/hashtags", srv.HashTagsHandler)
	profile.GET("/hashtags/:name", srv.HashTagHandler)
	profile.GET("/users", srv.UsersHandler)
	profile.GET("/users/:name", srv.UserHandler)
	profile.GET("/files", srv.FilesHandler)
	profile.GET("/notes", srv.NotesHandler)
	profile.GET("/archives", srv.ArchivesHandler)
	profile.GET("/archives/:date", srv.ArchiveNotesHandler)

	profile.GET("/emojis/:name", srv.EmojiHandler)

	return e
}

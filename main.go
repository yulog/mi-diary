package main

import (
	"context"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/infra"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/migrate"
	"github.com/yulog/mi-diary/server"
)

const name = "mi-diary"

const version = "0.0.1"

var revision = "HEAD"

func main() {
	app := app.New()
	infra := infra.New(app)
	logic := logic.New(infra)
	srv := server.New(logic)

	migrate.Do(infra)

	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	// e.Validator = &Validator{validator: validator.New()}

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
	profile.GET("/reactions/:name", srv.ReactionsHandler)
	profile.GET("/hashtags/:name", srv.HashTagsHandler)
	profile.GET("/users/:name", srv.UsersHandler)
	profile.GET("/files", srv.FilesHandler)
	profile.GET("/notes", srv.NotesHandler)
	profile.GET("/archives", srv.ArchivesHandler)
	profile.GET("/archives/:date", srv.ArchiveNotesHandler)
	profile.GET("/settings", srv.SettingsHandler)
	profile.POST("/settings/reactions", srv.SettingsReactionsHandler)
	profile.POST("/settings/emojis", srv.SettingsEmojisHandler)

	// TODO: context良く分からない
	go logic.JobProcesser(context.Background())

	e.Logger.Fatal(e.Start(net.JoinHostPort("", app.Config.Port)))
}

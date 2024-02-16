package main

import (
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

	profiles := e.Group("/profiles")
	profiles.GET("", srv.NewProfilesHandler)
	profiles.POST("", srv.AddProfileHandler)

	profile := profiles.Group("/:profile")
	profile.GET("", srv.HomeHandler)
	profile.GET("/reactions/:name", srv.ReactionsHandler)
	profile.GET("/hashtags/:name", srv.HashTagsHandler)
	profile.GET("/users/:name", srv.UsersHandler)
	profile.GET("/notes", srv.NotesHandler)
	profile.GET("/archives", srv.ArchivesHandler)
	profile.GET("/archives/:date", srv.ArchiveNotesHandler)
	profile.GET("/settings", srv.SettingsHandler)
	profile.POST("/settings/reactions", srv.SettingsReactionsHandler)
	profile.POST("/settings/emojis", srv.SettingsEmojisHandler)

	e.Logger.Fatal(e.Start(net.JoinHostPort("", app.Config.Port)))
}

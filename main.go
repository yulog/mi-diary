package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/migrate"
	"github.com/yulog/mi-diary/server"
)

const name = "mi-diary"

const version = "0.0.1"

var revision = "HEAD"

func main() {
	migrate.RunMigrations()
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	app := app.New()
	srv := server.New(app)
	// e.Validator = &Validator{validator: validator.New()}
	e.GET("/", srv.IndexHandler)
	e.GET("/reactions/:name", srv.ReactionsHandler)
	e.GET("/hashtags/:name", srv.HashTagsHandler)
	e.GET("/users/:name", srv.UsersHandler)
	e.GET("/notes", srv.NotesHandler)
	e.GET("/archives", srv.ArchivesHandler)
	e.GET("/archives/:date", srv.ArchiveNotesHandler)
	e.GET("/settings", srv.SettingsHandler)
	e.POST("/settings/reactions", srv.SettingsReactionsHandler)
	e.POST("/settings/emojis", srv.SettingsEmojisHandler)
	e.Logger.Fatal(e.Start(":" + app.Config.Port)) // TODO: configで変えられるようにする
}

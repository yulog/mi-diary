package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/migrate"
	"github.com/yulog/mi-diary/server"
)

const name = "mi-diary"

const version = "0.0.1"

var revision = "HEAD"

func main() {
	app := app.New()
	logic := logic.New(app)
	srv := server.New(logic)
	migrate.Do(app)
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	// e.Validator = &Validator{validator: validator.New()}
	e.GET("/", srv.ProfileHandler)
	e.GET("/:profile", srv.HomeHandler)
	e.GET("/:profile/reactions/:name", srv.ReactionsHandler)
	e.GET("/:profile/hashtags/:name", srv.HashTagsHandler)
	e.GET("/:profile/users/:name", srv.UsersHandler)
	e.GET("/:profile/notes", srv.NotesHandler)
	e.GET("/:profile/archives", srv.ArchivesHandler)
	e.GET("/:profile/archives/:date", srv.ArchiveNotesHandler)
	e.GET("/:profile/settings", srv.SettingsHandler)
	e.POST("/:profile/settings/reactions", srv.SettingsReactionsHandler)
	e.POST("/:profile/settings/emojis", srv.SettingsEmojisHandler)
	e.Logger.Fatal(e.Start(":" + app.Config.Port))
}

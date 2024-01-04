package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/server"
)

func main() {
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
	e.Logger.Fatal(e.Start(":1323"))
}

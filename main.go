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
	e.GET("/", srv.GetIndex)
	e.GET("/reactions/:name", srv.GetReactions)
	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"context"
	"net"

	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/infra"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/migrate"
	"github.com/yulog/mi-diary/server"
)

const name = "mi-diary"

const version = "0.0.3"

var revision = "HEAD"

func main() {
	app := app.New()
	infra := infra.New(app)
	logic := logic.New(infra)
	srv := server.New(logic)

	migrate.Do(infra)

	e := srv.NewRouter()

	// TODO: context良く分からない
	go logic.JobProcesser(context.Background())

	e.Logger.Fatal(e.Start(net.JoinHostPort("", app.Config.Port)))
}

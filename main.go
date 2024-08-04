package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/charmbracelet/log"
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
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger := slog.New(log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true}))
	slog.SetDefault(logger)
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

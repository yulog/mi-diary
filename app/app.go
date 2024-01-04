package app

import (
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/yulog/mi-diary/model"
)

type App struct {
	dbOnce sync.Once
	db     *bun.DB
}

func New() *App {
	return &App{}
}

func (app *App) DB() *bun.DB {
	app.dbOnce.Do(func() {
		// sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
		sqldb, err := sql.Open(sqliteshim.ShimName, "file:diary.db")
		if err != nil {
			panic(err)
		}
		db := bun.NewDB(sqldb, sqlitedialect.New())
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
		// modelを最初に使う前にやる
		db.RegisterModel((*model.NoteToTag)(nil))

		app.db = db
	})
	return app.db
}

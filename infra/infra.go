package infra

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/model"
)

type Infra struct {
	app *app.App

	db sync.Map // TODO:  sync.Onceの代わりになるのか？
}

func New(a *app.App) *Infra {
	return &Infra{app: a}
}

func (infra *Infra) Config() *app.Config {
	return &infra.app.Config
}

func (infra *Infra) DB(profile string) *bun.DB {
	v, _ := infra.db.LoadOrStore(profile, connect(profile))
	return v.(*bun.DB)
}

func connect(profile string) *bun.DB {
	// sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:diary_%s.db", profile))
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	// modelを最初に使う前にやる
	db.RegisterModel(
		(*model.NoteToTag)(nil),
		(*model.NoteToFile)(nil),
	)

	return db
}

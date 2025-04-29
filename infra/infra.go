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
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/migrate"
)

type Infra struct {
	app *app.App

	db sync.Map // TODO:  sync.Onceの代わりになるのか？
}

func New(a *app.App) logic.Repositorier {
	return &Infra{app: a}
}

func (infra *Infra) DB(profile string) *bun.DB {
	v, _ := infra.db.LoadOrStore(profile, connect(profile))
	return v.(*bun.DB)
}

func (infra *Infra) GenerateSchema(profile string) {
	migrate.GenerateSchema(infra.DB(profile))
}

func (infra *Infra) Migrate(profile string) {
	migrate.Do(infra.DB(profile).DB)
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

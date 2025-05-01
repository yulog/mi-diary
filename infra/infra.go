package infra

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/logic"
)

type Infra struct {
	app *app.App

	dao *DataBase
}

func New(a *app.App) logic.Repositorier {
	return &Infra{
		app: a,
		dao: NewDAO(),
	}
}

type DataBase struct {
	db sync.Map // TODO:  sync.Onceの代わりになるのか？
}

func NewDAO() *DataBase {
	return &DataBase{}
}

func (infra *DataBase) DB(profile string) *bun.DB {
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

func (i *Infra) NewUnitOfWorkInfra() repository.UnitOfWorkRepositorier {
	return i.dao
}

type txKey struct{}

func txFromContext(ctx context.Context) (bun.IDB, bool) {
	tx, ok := ctx.Value(txKey{}).(bun.IDB)
	return tx, ok
}

func (infra *DataBase) RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error) {
	err := infra.DB(profile).RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		ctx = context.WithValue(ctx, txKey{}, tx)
		if err := fn(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

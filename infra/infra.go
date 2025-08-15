package infra

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/akrylysov/pogreb"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/logic"
)

type Infra struct {
	dao *DataBase
}

func New() logic.Repositorier {
	return &Infra{
		dao: NewDAO(),
	}
}

type DataBase struct {
	db    sync.Map // TODO:  sync.Onceの代わりになるのか？
	cache sync.Map
}

func NewDAO() *DataBase {
	return &DataBase{}
}

func (dao *DataBase) DB(profile string) *bun.DB {
	v, _ := dao.db.LoadOrStore(profile, connect(profile))
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

func (dao *DataBase) Cache(host string) *pogreb.DB {
	v, ok := dao.cache.Load(host)
	if ok {
		return v.(*pogreb.DB)
	}
	db := connectCache(host)
	dao.cache.Store(host, db)

	return db
}

func connectCache(host string) *pogreb.DB {
	log.Println("connect")
	db, err := pogreb.Open(fmt.Sprintf(".cache/%s.db", host), nil)
	if err != nil {
		log.Println(err)
		panic(err)
	}

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

func (dao *DataBase) RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error) {
	err := dao.DB(profile).RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
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

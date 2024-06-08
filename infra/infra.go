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

// TODO: Config を直接返すのと個別に返すのとどちらが良い？
// func (infra *Infra) Config() *app.Config {
// 	return &infra.app.Config
// }

func (infra *Infra) SetConfig(key string, prof app.Profile) {
	infra.app.Config.Profiles[key] = prof
}

func (infra *Infra) StoreConfig() error {
	return app.ForceWriteConfig(&infra.app.Config)
}

func (infra *Infra) GetPort() string {
	return infra.app.Config.Port
}

func (infra *Infra) GetProfile(key string) (app.Profile, error) {
	v, ok := infra.app.Config.Profiles[key]
	if !ok {
		return app.Profile{}, fmt.Errorf("invalid profile: %s", key)
	}
	return v, nil
}

func (infra *Infra) GetProfileHost(key string) (string, error) {
	v, ok := infra.app.Config.Profiles[key]
	if !ok {
		return "", fmt.Errorf("invalid profile: %s", key)
	}
	return v.Host, nil
}

func (infra *Infra) GetProfiles() *app.Profiles {
	return &infra.app.Config.Profiles
}

func (infra *Infra) GetProgress() (int, int) {
	infra.app.Progress.RLock()
	defer infra.app.Progress.RUnlock()
	return infra.app.Progress.Progress, infra.app.Progress.Total
}

func (infra *Infra) SetProgress(p, t int) (int, int) {
	infra.app.Progress.Lock()
	defer infra.app.Progress.Unlock()
	infra.app.Progress.Progress = p
	infra.app.Progress.Total = t
	return p, t
}

func (infra *Infra) GetProgressDone() bool {
	infra.app.Progress.RLock()
	defer infra.app.Progress.RUnlock()
	return infra.app.Progress.Done
}

func (infra *Infra) SetProgressDone(d bool) bool {
	infra.app.Progress.Lock()
	defer infra.app.Progress.Unlock()
	infra.app.Progress.Done = d
	return d
}

func (infra *Infra) GetJob() chan app.Job {
	return infra.app.Job
}

func (infra *Infra) SetJob(j app.Job) {
	infra.app.Job <- j
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

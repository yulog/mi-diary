package app

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/yulog/mi-diary/model"
)

type App struct {
	Config Config

	dbOnce sync.Once
	db     *bun.DB
}

type Config struct {
	I      string
	UserId string
	Host   string
}

func New() *App {
	return &App{
		Config: LoadConfig(),
	}
}

func LoadConfig() Config {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("i", "")
	viper.SetDefault("userId", "")
	viper.SetDefault("host", "")
	viper.SafeWriteConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	return config
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

package app

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

type App struct {
	Config Config

	// move to infra
	// db sync.Map // TODO:  sync.Onceの代わりになるのか？
}

type Config struct {
	Port     string
	Profiles map[string]Profile
}

type Profile struct {
	I        string
	UserId   string
	UserName string
	Host     string
}

func New() *App {
	return &App{
		Config: loadConfig(),
	}
}

func loadConfig() Config {
	cfg := &Config{
		Port: "1323",
		Profiles: map[string]Profile{
			"default": {
				I:        "",
				UserId:   "",
				UserName: "",
				Host:     "",
			},
		},
	}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	b, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	viper.ReadConfig(bytes.NewBuffer(b))
	viper.SafeWriteConfig()
	err = viper.ReadInConfig()
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

// move to infra
// func (app *App) DB(profile string) *bun.DB {
// 	v, _ := app.db.LoadOrStore(profile, connect(profile))
// 	return v.(*bun.DB)
// }

// move to infra
// func connect(profile string) *bun.DB {
// 	// sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
// 	sqldb, err := sql.Open(sqliteshim.ShimName, fmt.Sprintf("file:diary_%s.db", profile))
// 	if err != nil {
// 		panic(err)
// 	}
// 	db := bun.NewDB(sqldb, sqlitedialect.New())
// 	db.AddQueryHook(bundebug.NewQueryHook(
// 		bundebug.WithVerbose(true),
// 		bundebug.FromEnv("BUNDEBUG"),
// 	))
// 	// modelを最初に使う前にやる
// 	db.RegisterModel((*model.NoteToTag)(nil))

// 	return db
// }

package app

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

type App struct {
	Config Config
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

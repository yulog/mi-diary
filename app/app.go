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
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	b, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	v.ReadConfig(bytes.NewBuffer(b))
	v.SafeWriteConfig()
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	return config
}

func ForceWriteConfig(cfg *Config) error {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	v.ReadConfig(bytes.NewBuffer(b))
	v.WriteConfig()
	return nil
}

package app

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

type App struct {
	Config   Config
	Job      chan Job
	Progress *Progress
}

type Config struct {
	Port     string
	Profiles Profiles
}

type Profiles map[string]Profile

type Profile struct {
	I        string
	UserId   string
	UserName string
	Host     string
}

type JobType int

const (
	Reaction JobType = iota + 1
	ReactionOne
	ReactionFull
	Emoji
)

func (j JobType) String() string {
	switch j {
	case Reaction:
		return "reaction"
	case ReactionOne:
		return "reaction(one)"
	case ReactionFull:
		return "reaction(full scan)"
	case Emoji:
		return "emoji"
	default:
		return "unkown"
	}
}

type Job struct {
	Profile string
	Type    JobType
	ID      string
}

type Progress struct {
	sync.RWMutex
	Progress int
	Total    int
	Done     bool
}

func New() *App {
	return &App{
		Config:   loadConfig(),
		Job:      make(chan Job),
		Progress: &Progress{},
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

func (cfg *Config) ForceWriteConfig() error {
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

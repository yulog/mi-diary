package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Profiles Profiles
}

type Profiles map[string]Profile

type Profile struct {
	I        string
	UserID   string
	UserName string
	Host     string
}

func (p Profiles) Get(key string) (Profile, error) {
	v, ok := p[key]
	if !ok {
		return Profile{}, fmt.Errorf("invalid profile: %s", key)
	}
	return v, nil
}

func DefaultConfig() *Config {
	return &Config{
		Port: "1323",
	}
}

func viperInstance() *viper.Viper {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	return v
}

func Load() *Config {
	v := viperInstance()

	b, err := json.Marshal(DefaultConfig())
	if err != nil {
		panic(err)
	}

	v.ReadConfig(bytes.NewBuffer(b))
	v.SafeWriteConfig() // ファイルがなかったときだけ書き込み
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	config := &Config{}
	err = v.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}
	return config
}

func (cfg *Config) Write() error {
	v := viperInstance()

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = v.ReadConfig(bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	err = v.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}

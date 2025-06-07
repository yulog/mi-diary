package app

import (
	"github.com/yulog/mi-diary/internal/config"
)

type App struct {
	Config config.Config
}

// type Config struct {
// 	Port     string
// 	Profiles Profiles
// }

// type Profiles map[string]Profile

// type Profile struct {
// 	I        string
// 	UserId   string
// 	UserName string
// 	Host     string
// }

// func (p Profiles) Get(key string) (Profile, error) {
// 	v, ok := p[key]
// 	if !ok {
// 		return Profile{}, fmt.Errorf("invalid profile: %s", key)
// 	}
// 	return v, nil
// }

func New() *App {
	return &App{
		Config: *config.Load(),
	}
}

// func loadConfig() Config {
// 	cfg := &Config{
// 		Port: "1323",
// 	}
// 	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
// 	v.SetConfigName("config")
// 	v.SetConfigType("json")
// 	v.AddConfigPath(".")

// 	b, err := json.Marshal(cfg)
// 	if err != nil {
// 		panic(err)
// 	}

// 	v.ReadConfig(bytes.NewBuffer(b))
// 	v.SafeWriteConfig()
// 	err = v.ReadInConfig()
// 	if err != nil {
// 		panic(fmt.Errorf("fatal error config file: %w", err))
// 	}
// 	var config Config
// 	err = v.Unmarshal(&config)
// 	if err != nil {
// 		panic(fmt.Errorf("unable to decode into struct, %v", err))
// 	}
// 	return config
// }

// func (cfg *Config) ForceWriteConfig() error {
// 	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
// 	v.SetConfigName("config")
// 	v.SetConfigType("json")
// 	v.AddConfigPath(".")

// 	b, err := json.Marshal(cfg)
// 	if err != nil {
// 		return err
// 	}

// 	v.ReadConfig(bytes.NewBuffer(b))
// 	v.WriteConfig()
// 	return nil
// }

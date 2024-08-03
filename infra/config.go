package infra

import (
	"fmt"

	"github.com/yulog/mi-diary/app"
)

// TODO: Config を直接返すのと個別に返すのとどちらが良い？
// func (infra *Infra) Config() *app.Config {
// 	return &infra.app.Config
// }

func (infra *Infra) SetConfig(key string, prof app.Profile) {
	infra.app.Config.Profiles[key] = prof
}

func (infra *Infra) StoreConfig() error {
	return infra.app.Config.ForceWriteConfig()
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
	p, err := infra.GetProfile(key)
	if err != nil {
		return "", err
	}
	return p.Host, nil
}

func (infra *Infra) GetProfiles() *app.Profiles {
	return &infra.app.Config.Profiles
}

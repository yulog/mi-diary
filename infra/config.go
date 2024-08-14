package infra

import (
	"fmt"
	"maps"
	"slices"

	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/logic"
)

type ConfigInfra struct {
	app *app.App
}

func NewConfigInfra(a *app.App) logic.ConfigRepositorier {
	return &ConfigInfra{app: a}
}

// TODO: Config を直接返すのと個別に返すのとどちらが良い？
// func (infra *Infra) Config() *app.Config {
// 	return &infra.app.Config
// }

func (infra *ConfigInfra) SetConfig(key string, prof app.Profile) {
	infra.app.Config.Profiles[key] = prof
}

func (infra *ConfigInfra) StoreConfig() error {
	return infra.app.Config.ForceWriteConfig()
}

func (infra *ConfigInfra) GetPort() string {
	return infra.app.Config.Port
}

func (infra *ConfigInfra) GetProfile(key string) (app.Profile, error) {
	v, ok := infra.app.Config.Profiles[key]
	if !ok {
		return app.Profile{}, fmt.Errorf("invalid profile: %s", key)
	}
	return v, nil
}

func (infra *ConfigInfra) GetProfileHost(key string) (string, error) {
	p, err := infra.GetProfile(key)
	if err != nil {
		return "", err
	}
	return p.Host, nil
}

func (infra *ConfigInfra) GetProfiles() *app.Profiles {
	return &infra.app.Config.Profiles
}

// range over funcでSortしたmapを使える
// 使う予定はない
func (infra *ConfigInfra) GetProfilesSorted(yield func(k string, v app.Profile) bool) {
	for _, k := range slices.Sorted(maps.Keys(infra.app.Config.Profiles)) {
		if !yield(k, infra.app.Config.Profiles[k]) {
			return
		}
	}
}

func (infra *ConfigInfra) GetProfilesSortedKey() []string {
	return slices.Sorted(maps.Keys(infra.app.Config.Profiles))
}

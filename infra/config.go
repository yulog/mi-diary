package infra

import (
	"maps"
	"slices"

	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/internal/config"
)

type ConfigInfra struct {
	app *app.App
}

func NewConfigInfra(a *app.App) repository.ConfigRepositorier {
	return &ConfigInfra{app: a}
}

// TODO: Config を直接返すのと個別に返すのとどちらが良い？
// func (infra *Infra) Config() *app.Config {
// 	return &infra.app.Config
// }

func (infra *ConfigInfra) SetConfig(key string, prof config.Profile) {
	infra.app.Config.Profiles[key] = prof
}

func (infra *ConfigInfra) StoreConfig() error {
	return infra.app.Config.Write()
}

func (infra *ConfigInfra) GetPort() string {
	return infra.app.Config.Port
}

func (infra *ConfigInfra) GetProfile(key string) (config.Profile, error) {
	return infra.app.Config.Profiles.Get(key)
}

func (infra *ConfigInfra) GetProfileHost(key string) (string, error) {
	p, err := infra.GetProfile(key)
	if err != nil {
		return "", err
	}
	return p.Host, nil
}

func (infra *ConfigInfra) GetProfiles() *config.Profiles {
	return &infra.app.Config.Profiles
}

// range over funcでSortしたmapを使える
// 使う予定はない
func (infra *ConfigInfra) GetProfilesSorted(yield func(k string, v config.Profile) bool) {
	for _, k := range slices.Sorted(maps.Keys(infra.app.Config.Profiles)) {
		if !yield(k, infra.app.Config.Profiles[k]) {
			return
		}
	}
}

func (infra *ConfigInfra) GetProfilesSortedKey() []string {
	return slices.Sorted(maps.Keys(infra.app.Config.Profiles))
}

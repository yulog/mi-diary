package repository

import "github.com/yulog/mi-diary/app"

type ConfigRepositorier interface {
	SetConfig(key string, prof app.Profile)
	StoreConfig() error

	GetProfiles() *app.Profiles
	GetProfilesSorted(yield func(k string, v app.Profile) bool)
	GetProfilesSortedKey() []string
	GetProfile(key string) (app.Profile, error)
	GetProfileHost(key string) (string, error)
	GetPort() string
}

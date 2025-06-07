package repository

import (
	"github.com/yulog/mi-diary/internal/config"
)

type ConfigRepositorier interface {
	SetConfig(key string, prof config.Profile)
	StoreConfig() error

	GetProfiles() *config.Profiles
	GetProfilesSorted(yield func(k string, v config.Profile) bool)
	GetProfilesSortedKey() []string
	GetProfile(key string) (config.Profile, error)
	GetProfileHost(key string) (string, error)
	GetPort() string
}

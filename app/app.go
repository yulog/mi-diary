package app

import (
	"github.com/yulog/mi-diary/internal/config"
)

type App struct {
	Config config.Config
}

func New() *App {
	return &App{
		Config: *config.Load(),
	}
}

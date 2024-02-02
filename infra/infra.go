package infra

import (
	"github.com/yulog/mi-diary/app"
)

type Infra struct {
	app *app.App
}

func New(a *app.App) *Infra {
	return &Infra{app: a}
}

func (infra *Infra) Config() *app.Config {
	return &infra.app.Config
}

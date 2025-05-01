package infra

import (
	"github.com/yulog/mi-diary/domain/service"
	"github.com/yulog/mi-diary/migrate"
)

type MigrationInfra struct {
	infra *DataBase
}

func (i *Infra) NewMigrationInfra() service.MigrationServicer {
	return &MigrationInfra{infra: i.dao}
}

func (i *MigrationInfra) GenerateSchema(profile string) {
	migrate.GenerateSchema(i.infra.DB(profile))
}

func (i *MigrationInfra) Execute(profile string) {
	migrate.Do(i.infra.DB(profile).DB)
}

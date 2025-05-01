package infra

import (
	"github.com/yulog/mi-diary/domain/service"
	"github.com/yulog/mi-diary/migrate"
)

type MigrationInfra struct {
	dao *DataBase
}

func (i *Infra) NewMigrationInfra() service.MigrationServicer {
	return &MigrationInfra{dao: i.dao}
}

func (i *MigrationInfra) GenerateSchema(profile string) {
	migrate.GenerateSchema(i.dao.DB(profile))
}

func (i *MigrationInfra) Execute(profile string) {
	migrate.Do(i.dao.DB(profile).DB)
}

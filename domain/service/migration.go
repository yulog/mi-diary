package service

type MigrationServicer interface {
	GenerateSchema(profile string)
	Execute(profile string)
}

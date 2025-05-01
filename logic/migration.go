package logic

func (l *Logic) GenerateSchema() {
	for k := range *l.ConfigRepo.GetProfiles() {
		l.MigrationService.GenerateSchema(k)
		break // schemaの生成は1つだけやれば良さそう
	}
}

func (l *Logic) Migrate() {
	for k := range *l.ConfigRepo.GetProfiles() {
		l.MigrationService.Execute(k)
	}
}

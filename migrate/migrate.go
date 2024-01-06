package migrate

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/yulog/mi-diary/app"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func RunMigrations() {
	app := app.New()
	migrations, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		panic(err)
	}

	driver, err := sqlite3.WithInstance(app.DB().DB, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}
	// m, _ := migrate.NewWithDatabaseInstance(
	//     "file:///migrations",
	//     "sqlite", driver)
	m, err := migrate.NewWithInstance(
		"iofs",
		migrations,
		"sqlite",
		driver,
	)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		panic(err)
	}
}

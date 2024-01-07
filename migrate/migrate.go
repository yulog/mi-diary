package migrate

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/yulog/mi-diary/app"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// RunMigrations はマイグレーションを実行する
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
		// 恐らく最新ということ
		log.Println(err)
		// panic(err)
	}
}

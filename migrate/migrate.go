package migrate

import (
	"database/sql"
	"embed"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Do はマイグレーションを実行する
func Do(db *sql.DB) {
	migrations, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		panic(err)
	}

	// 各プロファイルのDBをマイグレーションする
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
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
		slog.Info(err.Error())
	}
}

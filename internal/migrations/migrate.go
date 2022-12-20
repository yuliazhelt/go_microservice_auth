package migration

import (
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"log"
	"strings"
)

//go:embed *
var MmigrationsAssets embed.FS

func MigrateDb() {
	MigrationAssets, _ := fs.Sub(MmigrationsAssets, "migrations")

	sourceInstance, err := iofs.New(MigrationAssets, ".")
	if err != nil {
		log.Fatalf("cannot create source instance: %w", err)
	}

	db, err := sqlx.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	if err != nil {
		log.Fatalf("cannot open postgres: %w", err)
	}

	targetInstance, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("cannot target instance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceInstance, "postgres", targetInstance)
	if err != nil {
		log.Fatalf("cannot create migrate object: %w", err)
	}

	err = m.Up()
	if err != nil && !strings.Contains(err.Error(), "no change") {
		log.Fatalf("cannot migrate db: %w", err)
	}
}

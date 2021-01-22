package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbSource string, migrationsSource string) error {

	dbMigration, err := migrate.New(migrationsSource, dbSource)
	if err != nil {
		return fmt.Errorf("migrate.New error: %w", err)
	}

	err = dbMigration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate.Up error: %w", err)
	}

	return nil
}

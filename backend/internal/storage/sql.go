package storage

import (
	"embed"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type SQL struct {
	db *sqlx.DB
}

func NewSQL(db *sqlx.DB) SQL {
	return SQL{db: db}
}

//go:embed migrations/*.sql
var MigrationFiles embed.FS

type MigrationLogger struct{}

func (l MigrationLogger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l MigrationLogger) Verbose() bool {
	return true
}

// RunMigrations runs all migrations on startup
func (s SQL) RunMigrations() error {
	src, err := iofs.New(MigrationFiles, "migrations")
	if err != nil {
		return err
	}

	driver, err := pgx.WithInstance(s.db.DB, &pgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "pgx", driver)
	if err != nil {
		return err
	}

	m.Log = &MigrationLogger{}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DatabaseUtils interface {
	WaitForPostgresDb(host, port, databaseName string) (*sql.DB, error)
	RunMigrations(migrationsDir, host, port, databaseName string) error
	EnsureDatabaseExists(host, port, databaseName string) error
}

type DatabaseUtilsImpl struct{}

const (
	DefaultPostgresPort         = "5432"
	PostgresAdminDatabase       = "postgres"
	PostgresApplicationDatabase = "application"
	PostgresApplicationTemplate = "application_template"
)

func (d *DatabaseUtilsImpl) WaitForPostgresDb(host, port, databaseName string) (*sql.DB, error) {
	var err error
	var dbClient *sql.DB
	counter := 0
	attempts := 30

	for {
		counter++
		if counter >= attempts {
			return nil, Logger.NewError("failed to connect to database", "last_connections_error", err.Error())
		}

		dataSourceName := fmt.Sprintf("host=%s port=%s user=postgres dbname=%s sslmode=disable", host, port, databaseName)
		dbClient, err = sql.Open("postgres", dataSourceName)
		if err == nil {
			err = dbClient.Ping()
			if err == nil {
				break
			}
		}

		Logger.Info("waiting for postgres database to be ready", HostField, host, CurrentAttemptField, counter, MaximumAttemptsFields, attempts, PortField, port)
		time.Sleep(1 * time.Second)
	}

	return dbClient, nil
}

func (d *DatabaseUtilsImpl) RunMigrations(migrationsDir, host, port, databaseName string) error {
	migrator, err := migrate.New(
		"file://"+migrationsDir,
		fmt.Sprintf("postgres://postgres@%s:%s/%s?sslmode=disable", host, port, databaseName),
	)
	if err != nil {
		return Logger.NewError(err.Error())
	}
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return Logger.NewError(err.Error())
	}
	return nil
}

func (d *DatabaseUtilsImpl) EnsureDatabaseExists(host, port, databaseName string) error {
	adminDb, err := d.WaitForPostgresDb(host, port, "postgres")
	if err != nil {
		return err
	}
	defer Close(adminDb)

	var exists int
	existsError := adminDb.QueryRow(`SELECT 1 FROM pg_database WHERE datname = $1`, databaseName).Scan(&exists)
	if existsError == nil {
		return nil
	}

	_, createError := adminDb.Exec(`CREATE DATABASE ` + databaseName)
	if createError != nil {
		return Logger.NewError(createError.Error())
	}
	return nil
}

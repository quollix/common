package utils

import "database/sql"

type DatabaseConnector interface {
	GetDB() *sql.DB
	Connect() error
}

type DatabaseSnapshotRepository interface {
	CreateDatabaseSnapshot() error
	ResetDatabaseToSnapshot() error
	DeleteDatabaseSnapshotIfExists() error
	SnapshotExists() (bool, error)
}

type DatabaseSnapshotRepositoryImpl struct {
	DatabaseHost            string
	DatabasePort            string
	ApplicationDatabaseName string
	TemplateDatabaseName    string
	AdminDatabaseName       string
	DatabaseConnector       DatabaseConnector
	DatabaseUtils           DatabaseUtils
}

func (d *DatabaseSnapshotRepositoryImpl) CreateDatabaseSnapshot() error {
	adminDb, err := d.DatabaseUtils.WaitForPostgresDb(
		d.DatabaseHost,
		d.DatabasePort,
		d.AdminDatabaseName,
	)
	if err != nil {
		return err
	}
	defer Close(adminDb)

	terminateSQL := `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid()`

	_, err = adminDb.Exec(terminateSQL, d.TemplateDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	_, err = adminDb.Exec(`DROP DATABASE IF EXISTS ` + d.TemplateDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	_, err = adminDb.Exec(terminateSQL, d.ApplicationDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	_, err = adminDb.Exec(`CREATE DATABASE ` + d.TemplateDatabaseName + ` TEMPLATE ` + d.ApplicationDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	return nil
}

func NewDatabaseSnapshotRepository(
	databaseHost string,
	databaseConnector DatabaseConnector,
	databaseUtils DatabaseUtils,
) *DatabaseSnapshotRepositoryImpl {
	return &DatabaseSnapshotRepositoryImpl{
		DatabaseHost:            databaseHost,
		DatabasePort:            DefaultPostgresPort,
		ApplicationDatabaseName: PostgresApplicationDatabase,
		TemplateDatabaseName:    PostgresApplicationTemplate,
		AdminDatabaseName:       PostgresAdminDatabase,
		DatabaseConnector:       databaseConnector,
		DatabaseUtils:           databaseUtils,
	}
}

func (d *DatabaseSnapshotRepositoryImpl) ResetDatabaseToSnapshot() error {
	db := d.DatabaseConnector.GetDB()
	if db != nil {
		_ = db.Close()
	}

	adminDb, err := d.DatabaseUtils.WaitForPostgresDb(d.DatabaseHost, d.DatabasePort, d.AdminDatabaseName)
	if err != nil {
		return err
	}
	defer Close(adminDb)

	var templateExists int
	if adminDb.QueryRow(`SELECT 1 FROM pg_database WHERE datname = $1`, d.TemplateDatabaseName).Scan(&templateExists) != nil {
		return Logger.NewError("database snapshot does not exist", "template_db", d.TemplateDatabaseName)
	}

	terminateSQL := `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid()`

	_, err = adminDb.Exec(terminateSQL, d.ApplicationDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	_, err = adminDb.Exec(`DROP DATABASE IF EXISTS ` + d.ApplicationDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	_, err = adminDb.Exec(`CREATE DATABASE ` + d.ApplicationDatabaseName + ` TEMPLATE ` + d.TemplateDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	return d.DatabaseConnector.Connect()
}

func (d *DatabaseSnapshotRepositoryImpl) DeleteDatabaseSnapshotIfExists() error {
	adminDb, err := d.DatabaseUtils.WaitForPostgresDb(
		d.DatabaseHost,
		d.DatabasePort,
		d.AdminDatabaseName,
	)
	if err != nil {
		return err
	}
	defer Close(adminDb)

	terminateSQL := `SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid()`

	_, err = adminDb.Exec(terminateSQL, d.TemplateDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	_, err = adminDb.Exec(`DROP DATABASE IF EXISTS ` + d.TemplateDatabaseName)
	if err != nil {
		return Logger.NewError(err.Error())
	}

	return nil
}

func (d *DatabaseSnapshotRepositoryImpl) SnapshotExists() (bool, error) {
	adminDb, err := d.DatabaseUtils.WaitForPostgresDb(d.DatabaseHost, d.DatabasePort, d.AdminDatabaseName)
	if err != nil {
		return false, err
	}
	defer Close(adminDb)

	var databaseName string
	queryErr := adminDb.QueryRow(`SELECT datname FROM pg_database WHERE datname = $1`, d.TemplateDatabaseName).Scan(&databaseName)
	if queryErr != nil {
		return false, nil
	}
	return true, nil
}

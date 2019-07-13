package db

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"time"
	"ygp/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

// Database holds all required information about database as well as methods
type Database struct {
	db *sql.DB
}

var (
	DB Database
)

// Setup creates database connection
func Setup() *Database {
	DB = Database{}

	// open connection
	db, err := sql.Open("postgres", config.Cfg.DatabaseConnectionString)
	if err != nil {
		logrus.WithError(err).Fatalf("[DB] Can't connect to DB")
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	DB.db = db

	logrus.Printf("[DB] Connected to %s", config.Cfg.DatabaseConnectionString)
	return &DB
}

// Teardown shutdown all remaining connections
func Teardown() {
	err := DB.db.Close()
	if err != nil {
		logrus.WithError(err).Print("[DB] can't teardown")
	}
}

// SaveChannel insert row into DB
func (d Database) SaveChannel(ctx context.Context, channelID string, e error) error {
	now := time.Now()
	id, _ := uuid.NewV4()
	_, err := d.db.ExecContext(ctx, `INSERT INTO channels_logs (id, channel_id, date, error) VALUES ($1, $2, $3, $4)`,
		id, channelID, now.UTC(), e)
	if err != nil {
		logrus.WithError(err).Printf("[DB] On save channel")
		return err
	}
	logrus.Printf("[DB] Saved to db %s", channelID)
	return nil
}

// Migrate starts process of migrating sql files in order
func Migrate() error {
	driver, err := postgres.WithInstance(DB.db, &postgres.Config{})
	if err != nil {
		logrus.WithError(err).Printf("[DB] Migration can't start...")
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		logrus.WithError(err).Printf("[DB] Can't migrate to existing database...")
		return err
	}
	err = m.Steps(2)
	if err != nil {
		logrus.WithError(err).Printf("[DB] Migration can't start...")
		return err
	}
	return nil
}
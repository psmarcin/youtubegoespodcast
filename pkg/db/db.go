package db

import (
	"context"
	"database/sql"
	"time"
	"ytg/pkg/config"

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
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	DB.db = db

	return &DB
}

// SaveChannel insert row into DB
func (d Database) SaveChannel(ctx context.Context, channelID string, e error) error {
	now := time.Now()
	uuid, _ := uuid.NewV4()
	_, err := d.db.ExecContext(ctx, `INSERT INTO channels_logs (id, channel_id, date, error) VALUES ($1, $2, $3, $4)`,
		uuid, channelID, now.UTC(), e)
	if err != nil {
		logrus.WithError(err).Printf("[DB] On save channel")
		return err
	}
	logrus.Printf("[DB] Saved to db %s", channelID)
	return nil
}

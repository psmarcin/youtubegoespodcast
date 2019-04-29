package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Setup set default values for logrus
func Setup() {
	if os.Getenv("APP_ENV") == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
}

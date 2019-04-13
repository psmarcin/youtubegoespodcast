package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// Setup set default values for logrus
func Setup() {
	if os.Getenv("APP_ENV") == "production" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
}

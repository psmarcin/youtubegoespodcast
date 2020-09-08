package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Setup set default values for logrus
func Setup() {
	if os.Getenv("APP_ENV") == "production" {
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "time",
				log.FieldKeyLevel: "severity",
				log.FieldKeyMsg:   "message",
			},
		})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
}

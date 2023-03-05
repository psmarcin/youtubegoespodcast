package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Setup set default values for logrus
func Setup() {
	loggerType := os.Getenv("LOGGER_TYPE")

	switch loggerType {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
		return
	default:
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "time",
				log.FieldKeyLevel: "severity",
				log.FieldKeyMsg:   "message",
			},
		})
		return
	}
}

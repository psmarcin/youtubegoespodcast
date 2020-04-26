package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config stores all data from env variables
type Config struct {
	AppEnv              string   `env:"APP_ENV,required" envDefault:"development"`
	GoogleAPIKey        string   `env:"GOOGLE_API_KEY,required"`
	Port                string   `env:"PORT" envDefault:"3000"`
	FirestoreCollection string   `env:"FIRESTORE_COLLECTION" envDefault:"dev"`
	CorsOrigins         []string `env:"CORS_ORIGINS" envDefault:"http://localhost:8080" envSeparator:","`
}

// Cfg stores config values
var Cfg = Config{}

// Init load variable from .env file or environment variables
func Init() {
	_ = godotenv.Load(".env")
	logrus.Info("[CONFIG] Environment variables loaded")

	loadConfig()
}

func loadConfig() {
	if err := env.Parse(&Cfg); err != nil {
		logrus.WithError(err).Fatal("[CONFIG] Can't parse config values")
	}
}

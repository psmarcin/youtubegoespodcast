package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config stores all data from env variables
type Config struct {
	ApiRouterLoggerFormat string   `env:"ROUTER_LOGGER_FORMAT" envDefault:"{\"severity\": \"INFO\", \"time\": \"${time}\", \"level\":\"info\", \"source\": \"router\", \"path\": \"${route}\", \"url\": \"${url}\", \"method\": \"${method}\", \"latency\": \"${latency}\"}\n"`
	AppEnv                string   `env:"APP_ENV,required" envDefault:"development"`
	CorsOrigins           []string `env:"CORS_ORIGINS" envDefault:"http://localhost:8080" envSeparator:","`
	FirestoreCollection   string   `env:"FIRESTORE_COLLECTION" envDefault:"dev"`
	GoogleAPIKey          string   `env:"GOOGLE_API_KEY,required"`
	Port                  string   `env:"PORT" envDefault:"3000"`
}

// Cfg stores config values
var Cfg = Config{}
var l = logrus.WithField("source", "config")

// Init load variable from .env file or environment variables
func Init() {
	_ = godotenv.Load(".env")
	l.Info("environment variables loaded")

	loadConfig()
}

func loadConfig() {
	if err := env.Parse(&Cfg); err != nil {
		l.WithError(err).Fatal("can't parse config values")
	}
}

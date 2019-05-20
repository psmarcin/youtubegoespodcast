package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config stores all data from env variables
type Config struct {
	AppEnv       string `env:"APP_ENV,required" envDefault:"development"`
	GoogleAPIKey string `env:"GOOGLE_API_KEY,required"`
	Port         string `env:"PORT" envDefault:"3000"`
	RedisPort    string `env:"REDIS_PORT" envDefault:"6379"`
	RedisHost    string `env:"REDIS_HOST" envDefault:"localhost"`
	// IsProduction bool          `env:"PRODUCTION"`
	// Hosts        []string      `env:"HOSTS" envSeparator:":"`
	// Duration     time.Duration `env:"DURATION"`
	// TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
}

// Cfg stores config values
var Cfg = Config{}

// Init load variable from .env file or environment variables
func Init() {
	_ = godotenv.Load(".env")
	logrus.Info("[CONFIG] env variables loaded")
	loadConfig()
}

func loadConfig() {
	if err := env.Parse(&Cfg); err != nil {
		logrus.WithError(err).Fatal("[CONFIG] can't parse config values")
	}
}

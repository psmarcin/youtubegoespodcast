package main

import (
	"github.com/sirupsen/logrus"
	"ygp/pkg/api"
	"ygp/pkg/cache"
	"ygp/pkg/config"
	"ygp/pkg/logger"
)

func main() {
	// Config
	config.Init()
	// Logger
	logger.Setup()
	// Cache
	_, _ = cache.Connect()
	// API
	app := api.Start()

	logrus.Fatal(app.Listen(config.Cfg.Port))
}

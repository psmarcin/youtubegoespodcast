package main

import (
	"github.com/sirupsen/logrus"
	"github.com/psmarcin/youtubegoespodcast/pkg/api"
	"github.com/psmarcin/youtubegoespodcast/pkg/cache"
	"github.com/psmarcin/youtubegoespodcast/pkg/config"
	"github.com/psmarcin/youtubegoespodcast/pkg/logger"
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

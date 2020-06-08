package main

import (
	"github.com/psmarcin/youtubegoespodcast/pkg/api"
	"github.com/psmarcin/youtubegoespodcast/pkg/cache"
	"github.com/psmarcin/youtubegoespodcast/pkg/config"
	"github.com/psmarcin/youtubegoespodcast/pkg/logger"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/sirupsen/logrus"
)

var l = logrus.WithField("source", "cmd")

func main() {
	// Config
	config.Init()
	// Logger
	logger.Setup()
	// Cache
	_, _ = cache.Connect()
	// YouTube API
	_, err := youtube.New()
	if err != nil {
		l.WithError(err).Fatalf("can't connect to youtube service")
	}
	// API
	app := api.Start()

	logrus.Fatal(app.Listen(config.Cfg.Port))
}

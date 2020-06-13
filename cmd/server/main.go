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
	c, err := cache.Connect()
	if err != nil {
		l.WithError(err).Fatalf("can't connect to youtube service")
	}
	// YouTube API
	yt, err := youtube.New()
	if err != nil {
		l.WithError(err).Fatalf("can't connect to youtube service")
	}
	// dependencies
	deps := api.Dependencies{
		Cache:   c,
		YouTube: yt,
	}
	// API
	app := api.Start(deps)

	logrus.Fatal(app.Listen(config.Cfg.Port))
}

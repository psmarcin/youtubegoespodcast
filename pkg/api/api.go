package api

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"github.com/gofiber/requestid"
	"net/http"
	"time"

	"github.com/gofiber/fiber"
	"github.com/sirupsen/logrus"
	"ygp/pkg/config"
)

var l = logrus.WithField("source", "API")

// Start starts API
func Start() *fiber.App {
	appConfig := fiber.Settings{
		CaseSensitive: true,
		Immutable:     true,
		ReadTimeout:   1 * time.Second,
		WriteTimeout:  1 * time.Second,
		IdleTimeout:   1 * time.Second,
	}
	corsConfig := cors.Config{
		Filter:       nil,
		AllowOrigins: config.Cfg.CorsOrigins,
		AllowMethods: []string{"GET"},
	}
	logConfig := logger.Config{
		Format: config.Cfg.ApiRouterLoggerFormat,
		TimeFormat: time.RFC3339,
	}

	// init fiber application
	app := fiber.New(&appConfig)

	// middlewares
	app.Use(cors.New(corsConfig))
	app.Use(logger.New(logConfig))
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())

	// define routes
	app.Get("/", rootHandler)
	app.Get("/trending", TrendingHandler)
	app.Get("/channels", ChannelsHandler)
	app.Get("/video/:videoId/track.mp3", VideoHandler)
	app.Head("/video/:videoId/track.mp3", VideoHandler)
	app.Get("/feed/channel/:"+ParamChannelId, FeedHandler)
	app.Head("/feed/channel/:"+ParamChannelId, FeedHandler)

	// not found handler
	// Last middleware to match anything
	app.Use(func(c *fiber.Ctx) {
		c.SendStatus(http.StatusNotFound) // => 404 "Not Found"
	})

	l.WithField("port", config.Cfg.Port).Infof("started")
	return app
}

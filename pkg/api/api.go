package api

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/helmet"
	"github.com/gofiber/recover"
	"github.com/gofiber/requestid"
	"github.com/psmarcin/logger"
	"time"

	"github.com/gofiber/fiber"
	"github.com/psmarcin/youtubegoespodcast/pkg/config"
	"github.com/sirupsen/logrus"
)

var l = logrus.WithField("source", "API")

// Start creates and starts HTTP server
func Start() *fiber.App {
	serverHTTP := CreateHTTPServer()
	// define routes
	serverHTTP.Get("/", rootHandler)
	serverHTTP.Post("/", rootHandler)
	serverHTTP.Static("/assets", "./web/static")

	videoGroup := serverHTTP.Group("/video")
	videoGroup.Get("/:videoId/track.mp3", videoHandler)
	videoGroup.Head("/:videoId/track.mp3", videoHandler)

	feedGroup := serverHTTP.Group("/feed/channel")
	feedGroup.Get("/:"+ParamChannelId, feedHandler)
	feedGroup.Head("/:"+ParamChannelId, feedHandler)

	// error found handler
	serverHTTP.Use(errorHandler)

	l.WithField("port", config.Cfg.Port).Infof("started")
	return serverHTTP
}

// CreateHTTPServer creates configured HTTP server
func CreateHTTPServer() *fiber.App {
	appConfig := fiber.Settings{
		CaseSensitive: true,
		Immutable:     false,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  3 * time.Second,
		IdleTimeout:   1 * time.Second,
	}
	corsConfig := cors.Config{
		Filter:       nil,
		AllowOrigins: config.Cfg.CorsOrigins,
		AllowMethods: []string{"GET"},
	}
	logConfig := logger.Config{
		Format:     config.Cfg.ApiRouterLoggerFormat,
		TimeFormat: time.RFC3339,
	}

	// init fiber application
	serverHTTP := fiber.New(&appConfig)

	// middleware
	serverHTTP.Use(cors.New(corsConfig))
	serverHTTP.Use(logger.New(logConfig))
	serverHTTP.Use(recover.New())
	serverHTTP.Use(requestid.New())
	serverHTTP.Use(helmet.New())

	return serverHTTP
}

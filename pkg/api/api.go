package api

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"github.com/gofiber/requestid"
	"github.com/gofiber/template/html"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/psmarcin/youtubegoespodcast/pkg/feed"
	"time"

	"github.com/gofiber/fiber"
	"github.com/psmarcin/youtubegoespodcast/pkg/config"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	YouTube app.YouTubeService
	YTDL    feed.YTDLDependencies
}

var l = logrus.WithField("source", "API")

// Start creates and starts HTTP server
func Start(deps Dependencies) *fiber.App {
	serverHTTP := CreateHTTPServer()
	feedDeps := CreateDependencies(deps)
	rootDeps := CreateRootDependencies(deps)
	videoDependencies := CreateVideoDependencies(deps)

	// define routes
	serverHTTP.Get("/", rootHandler(rootDeps))
	serverHTTP.Post("/", rootHandler(rootDeps))

	videoGroup := serverHTTP.Group("/video")
	videoGroup.Get("/:videoId/track.mp3", videoHandler(videoDependencies))
	videoGroup.Head("/:videoId/track.mp3", videoHandler(videoDependencies))

	feedGroup := serverHTTP.Group("/feed/channel")
	feedGroup.Get("/:"+ParamChannelId, feedHandler(feedDeps))
	feedGroup.Head("/:"+ParamChannelId, feedHandler(feedDeps))

	// error found handler
	serverHTTP.Use(errorHandler)

	l.WithField("port", config.Cfg.Port).Infof("started")

	return serverHTTP
}

// CreateHTTPServer creates configured HTTP server
func CreateHTTPServer() *fiber.App {
	templateEngine := html.New("./web/templates", ".tmpl")

	appConfig := fiber.Settings{
		CaseSensitive: true,
		Immutable:     false,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  3 * time.Second,
		IdleTimeout:   1 * time.Second,
		Views:         templateEngine,
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

	serverHTTP.Static("/assets", "./web/static")

	return serverHTTP
}

func CreateDependencies(deps Dependencies) feed.Dependencies {
	fd := feed.Dependencies{
		YouTube: deps.YouTube,
		YTDL:    deps.YTDL,
	}
	return fd
}

func CreateRootDependencies(deps Dependencies) rootDependencies {
	return deps.YouTube
}

func CreateVideoDependencies(deps Dependencies) videoDependencies {
	return deps.YTDL
}

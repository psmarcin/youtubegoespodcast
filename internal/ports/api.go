package ports

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
	"github.com/psmarcin/youtubegoespodcast/internal/config"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	server         *fiber.App
	youTubeService app.YouTubeService
	ytdlService    feed.YTDLDependencies
}

func NewHttpServer(
	server *fiber.App,
	youTubeService app.YouTubeService,
	ytdlService feed.YTDLDependencies,
) HttpServer {
	return HttpServer{
		server,
		youTubeService,
		ytdlService,
	}
}

func (h HttpServer) Serve() *fiber.App {
	feedDeps := h.getFeedDependencies()
	rootDeps := h.getRootDependencies()
	videoDeps := h.getVideoDependencies()

	// define routes
	h.server.Get("/", rootHandler(rootDeps))
	h.server.Post("/", rootHandler(rootDeps))

	videoGroup := h.server.Group("/video")
	videoGroup.Get("/:videoId/track.mp3", videoHandler(videoDeps))
	videoGroup.Head("/:videoId/track.mp3", videoHandler(videoDeps))

	feedGroup := h.server.Group("/feed/channel")
	feedGroup.Get("/:"+ParamChannelId, feedHandler(feedDeps))
	feedGroup.Head("/:"+ParamChannelId, feedHandler(feedDeps))

	// error found handler
	h.server.Use(errorHandler)

	//l.WithField("port", config.Cfg.Port).Infof("started")
	return h.server
}

func (h HttpServer) getFeedDependencies() feed.Dependencies {
	return feed.Dependencies{
		YouTube: h.youTubeService,
		YTDL:    h.ytdlService,
	}
}

func (h HttpServer) getRootDependencies() rootDependencies {
	return h.youTubeService
}

func (h HttpServer) getVideoDependencies() videoDependencies {
	return h.ytdlService
}

type Dependencies struct {
	YouTube app.YouTubeService
	YTDL    feed.YTDLDependencies
}

var l = logrus.WithField("source", "API")

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

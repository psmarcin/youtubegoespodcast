package ports

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/config"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	server         *fiber.App
	youTubeService app.YouTubeService
	ytdlService    app.YTDLDependencies
}

func NewHttpServer(
	server *fiber.App,
	youTubeService app.YouTubeService,
	ytdlService app.YTDLDependencies,
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

func (h HttpServer) getFeedDependencies() app.FeedService {
	return app.NewFeedService(h.youTubeService, h.ytdlService)
}

func (h HttpServer) getRootDependencies() rootDependencies {
	return h.youTubeService
}

func (h HttpServer) getVideoDependencies() videoDependencies {
	return h.ytdlService
}

type Dependencies struct {
	YouTube app.YouTubeService
	YTDL    app.YTDLDependencies
}

var l = logrus.WithField("source", "API")

// CreateHTTPServer creates configured HTTP server
func CreateHTTPServer() *fiber.App {
	templateEngine := html.New("./web/templates", ".tmpl")

	appConfig := fiber.Config{
		CaseSensitive: true,
		Immutable:     false,
		Prefork:       false,
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  3 * time.Second,
		IdleTimeout:   1 * time.Second,
		Views:         templateEngine,
	}
	logConfig := logger.Config{
		Format:     config.Cfg.ApiRouterLoggerFormat,
		TimeFormat: time.RFC3339,
	}

	// init fiber application
	serverHTTP := fiber.New(appConfig)

	// middleware
	serverHTTP.Use(RequestContext())
	serverHTTP.Use(logger.New(logConfig))
	serverHTTP.Use(recover.New())
	serverHTTP.Use(requestid.New())
	serverHTTP.Use(helmet.New())

	serverHTTP.Static("/assets", "./web/static")

	return serverHTTP
}

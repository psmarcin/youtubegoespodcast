package ports

import (
	"embed"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	hel "github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/html"
	fiberOtel "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	fib "github.com/gofiber/fiber/v2"
	"github.com/psmarcin/youtubegoespodcast/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	ReadTimeout  = 5
	WriteTimeout = 3
	IdleTimeout  = 1
)

type HTTPServer struct {
	server         *fib.App
	youTubeService app.YouTubeService
	fileService    app.FileService
}

func NewHTTPServer(
	server *fib.App,
	youTubeService app.YouTubeService,
	fileService app.FileService,
) HTTPServer {
	return HTTPServer{
		server,
		youTubeService,
		fileService,
	}
}

func (h HTTPServer) Serve() *fib.App {
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
	feedGroup.Get("/:"+ParamChannelID, feedHandler(feedDeps))
	feedGroup.Head("/:"+ParamChannelID, feedHandler(feedDeps))

	// error found handler
	h.server.Use(errorHandler)

	// l.WithField("port", config.Cfg.Port).Infof("started")
	return h.server
}

func (h HTTPServer) getFeedDependencies() app.FeedService {
	return app.NewFeedService(h.youTubeService, h.fileService)
}

func (h HTTPServer) getRootDependencies() rootDependencies {
	return h.youTubeService
}

func (h HTTPServer) getVideoDependencies() videoDependencies {
	return h.fileService
}

type Dependencies struct {
	YouTube app.YouTubeService
	YTDL    app.YTDLDependencies
}

//go:embed templates/*
var templatesFs embed.FS

//go:embed static/*
var staticFS embed.FS
var l = logrus.WithField("source", "API")

// CreateHTTPServer creates configured HTTP server
func CreateHTTPServer() *fib.App {
	templateEngine := html.NewFileSystem(http.FS(templatesFs), ".tmpl")

	appConfig := fib.Config{
		CaseSensitive: true,
		Immutable:     false,
		Prefork:       false,
		ReadTimeout:   ReadTimeout * time.Second,
		WriteTimeout:  WriteTimeout * time.Second,
		IdleTimeout:   IdleTimeout * time.Second,
		Views:         templateEngine,
	}
	logConfig := logger.Config{
		Format:     config.Cfg.APIRouterLoggerFormat,
		TimeFormat: time.RFC3339,
	}

	// init fiber application
	serverHTTP := fib.New(appConfig)

	// middleware
	serverHTTP.Use(fiberOtel.New(fiberOtel.Config{
		Tracer: otel.GetTracerProvider().Tracer(
			"yt.psmarcin.dev/api",
			trace.WithInstrumentationVersion("1.0.0"),
		),
	}))
	serverHTTP.Use(logger.New(logConfig))
	serverHTTP.Use(recover.New())
	serverHTTP.Use(requestid.New())
	serverHTTP.Use(hel.New())

	serverHTTP.Use("/assets", filesystem.New(filesystem.Config{
		Root: http.FS(staticFS),
	}))

	return serverHTTP
}

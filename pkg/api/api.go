package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"ygp/pkg/config"
	"ygp/pkg/logger"
	"ygp/pkg/metrics"
)

const RootJSON = "{\"status\": \"OK\"}"

// Start starts API
func Start() {
	startMultiplex()
	logrus.Infof("[API] started on port :%s", config.Cfg.Port)
}

func startMultiplex() {
	// router init
	router := mux.NewRouter()

	// middleware
	router.Use(logger.Middleware)
	router.Use(MiddlewareCORS)
	router.Use(MiddlewareJSON)
	router.Use(metrics.MiddlewareHandlerCounter)
	router.Use(metrics.MiddlewareHandlerInFlight)
	router.Use(metrics.MiddlewareResponseSize)
	// routes
	router.Handle("/", metrics.MiddlewareHandlerDuration("/", rootHandler)).Methods("GET")
	router.Handle("/trending", metrics.MiddlewareHandlerDuration("/trending", TrendingHandler)).Methods("GET")
	router.Handle("/channels", metrics.MiddlewareHandlerDuration("/channels", ChannelsHandler)).Methods("GET")
	router.Handle("/video/{videoId}.mp3", metrics.MiddlewareHandlerDuration("/video/{videoId}.mp3", VideoHandlerMP3)).Methods("GET", "HEAD")
	router.Handle("/video/{videoId}.mp4", metrics.MiddlewareHandlerDuration("/video/{videoId}.mp4", VideoHandler)).Methods("GET", "HEAD")
	router.Handle("/video", metrics.MiddlewareHandlerDuration("/video", VideoRedirectHandler)).Methods("GET", "HEAD")
	router.Handle("/feed/channel/{channelId}", metrics.MiddlewareHandlerDuration("/feed/channel/{channelId}", FeedHandler)).Methods("GET", "HEAD")
	router.Handle("/metrics", promhttp.Handler())
	http.Handle("/", router)

	logrus.Infof("[API] Port %s", config.Cfg.Port)
	logrus.Fatal(http.ListenAndServe(":"+config.Cfg.Port, nil))
}

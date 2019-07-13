package api

import (
	"net/http"
	"ygp/pkg/channels"
	"ygp/pkg/config"
	"ygp/pkg/feed"
	"ygp/pkg/logger"
	"ygp/pkg/metrics"
	"ygp/pkg/trending"
	"ygp/pkg/utils"
	"ygp/pkg/video"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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
	router.Use(utils.MiddlewareCORS)
	router.Use(utils.MiddlewareJSON)
	router.Use(metrics.MiddlewareHandlerCounter)
	router.Use(metrics.MiddlewareHandlerInFlight)
	router.Use(metrics.MiddlewareResponseSize)
	// routes
	router.Handle("/", metrics.MiddlewareHandlerDuration("/", rootHandler)).Methods("GET")
	router.Handle("/trending", metrics.MiddlewareHandlerDuration("/trending", trending.Handler)).Methods("GET")
	router.Handle("/channels", metrics.MiddlewareHandlerDuration("/channels", channels.Handler)).Methods("GET")
	router.Handle("/video/{videoId}.mp3", metrics.MiddlewareHandlerDuration("/video/{videoId}.mp3", video.Handler)).Methods("GET", "HEAD")
	router.Handle("/video/{videoId}.mp4", metrics.MiddlewareHandlerDuration("/video/{videoId}.mp4", video.Handler)).Methods("GET", "HEAD")
	router.Handle("/video", metrics.MiddlewareHandlerDuration("/video", video.RedirectHandler)).Methods("GET", "HEAD")
	router.Handle("/feed/channel/{channelId}", metrics.MiddlewareHandlerDuration("/feed/channel/{channelId}", feed.Handler)).Methods("GET", "HEAD")
	router.Handle("/metrics", promhttp.Handler())
	http.Handle("/", router)
	logrus.Infof("[API] Port %s", config.Cfg.Port)
	logrus.Fatal(http.ListenAndServe(":"+config.Cfg.Port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w)
	utils.Send(w, RootJSON, http.StatusOK)
}

package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"ytg/pkg/channels"
	"ytg/pkg/config"
	"ytg/pkg/logger"
	"ytg/pkg/trending"
	"ytg/pkg/utils"
	"ytg/pkg/video"
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
	router.Use(utils.MiddlewareJSON)
	router.Use(utils.MiddlewareCORS)
	// routes
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/trending", trending.Handler).Methods("GET")
	router.HandleFunc("/channels", channels.Handler).Methods("GET")
	router.HandleFunc("/video/{videoId}.mp3", video.Handler).Methods("GET")
	router.HandleFunc("/video", video.RedirectHandler).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())
	http.Handle("/", router)
	logrus.Infof("[API] Port %s", config.Cfg.Port)
	logrus.Fatal(http.ListenAndServe(":"+config.Cfg.Port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w)
	utils.OkResponse(w)
	utils.WriteBodyResponse(w, RootJSON)
}

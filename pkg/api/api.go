package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"ytg/pkg/channels"
	"ytg/pkg/config"
	"ytg/pkg/trending"
	"ytg/pkg/utils"
)

const RootJSON = "{\"status\": \"OK\"}"

// Start starts API
func Start() {
	startMultiplex()
	logrus.Infof("[API] started on port :%s", config.Cfg.Port)
}

func startMultiplex() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/trending", trending.Handler)
	http.HandleFunc("/channels", channels.Handler)
	logrus.Infof("[API] Port %s", config.Cfg.Port)
	logrus.Fatal(http.ListenAndServe(":"+config.Cfg.Port, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("[API] request %s %s", r.Method, r.RequestURI)
	utils.JSONResponse(w)
	utils.OkResponse(w)
	utils.WriteBodyResponse(w, RootJSON)
}

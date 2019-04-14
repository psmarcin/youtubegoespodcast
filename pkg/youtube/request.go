package youtube

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"ytg/pkg/config"
)

const YouTubeURL = "https://www.googleapis.com/youtube/v3/"

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// Request do request to youtube server to get trending videos
func Request(req *http.Request, y interface{}) {
	logrus.Infof("[YT] Request to %s", req.URL.String())

	// Set API key
	query := req.URL.Query()
	query.Add("key", config.Cfg.GoogleAPIKey)
	req.URL.RawQuery = query.Encode()

	// Do request
	res, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("[YT] while doing request to youtube api")
	}
	defer res.Body.Close()

	// Decode
	err = json.NewDecoder(res.Body).Decode(y)
	if err != nil {
		logrus.WithError(err).Warn("[YT] while doing parsing body content")
	}
}

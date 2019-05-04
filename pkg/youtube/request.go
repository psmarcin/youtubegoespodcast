package youtube

import (
	"encoding/json"
	"net/http"
	"time"
	"ytg/pkg/config"

	"github.com/sirupsen/logrus"
)

const YouTubeURL = "https://www.googleapis.com/youtube/v3/"

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// Request do request to youtube server to get trending videos
func Request(req *http.Request, y interface{}) error {
	// Set API key
	query := req.URL.Query()
	query.Add("key", config.Cfg.GoogleAPIKey)
	req.URL.RawQuery = query.Encode()
	// Do request
	res, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("[YT] while doing request to youtube api")
		return err
	}
	defer res.Body.Close()

	// Decode
	err = json.NewDecoder(res.Body).Decode(y)
	if err != nil {
		logrus.WithError(err).Warn("[YT] while doing parsing body content")
		return err
	}
	return nil
}

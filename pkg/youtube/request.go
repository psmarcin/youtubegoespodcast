package youtube

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"ygp/pkg/config"
	"ygp/pkg/errx"

	"github.com/sirupsen/logrus"
)

const YouTubeURL = "https://www.googleapis.com/youtube/v3/"

var client = &http.Client{
	Timeout: 5 * time.Second,
}

// Request do request to youtube server to get trending videos
func Request(req *http.Request, y interface{}) errx.APIError {
	// Set API key
	query := req.URL.Query()
	query.Add("key", config.Cfg.GoogleAPIKey)
	req.URL.RawQuery = query.Encode()
	// Do request
	res, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Warn("[YT] while doing request to youtube api")
		return errx.NewAPIError(err, http.StatusInternalServerError)
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		logrus.WithError(err).Printf("[YT] Request error and response: %+v", res)
		return errx.NewAPIError(err, res.StatusCode)
	}
	defer res.Body.Close()

	// Decode
	err = json.NewDecoder(res.Body).Decode(y)
	if err != nil {
		logrus.WithError(err).Warn("[YT] while doing parsing body content")
		return errx.NewAPIError(err, http.StatusInternalServerError)
	}
	return errx.APIError{}
}

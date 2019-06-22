package youtube

import (
	"encoding/json"
	"net/http"
	"time"
	"ytg/pkg/errx"
	"ytg/pkg/redis_client"

	"github.com/sirupsen/logrus"
)

const (
	trendingCacheTTL    = time.Hour * 24
	trendingCachePrefix = "yttrending_main"
)

// GetTrendings collects trendings from YouTube API
func GetTrendings() (YoutubeResponse, errx.APIError) {
	trendings := YoutubeResponse{}
	_, err := redis_client.Client.GetKey(trendingCachePrefix, &trendings)
	if err == nil {
		return trendings, errx.APIError{}
	}

	req, err := http.NewRequest("GET", YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trendings, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("chart", "mostPopular")
	req.URL.RawQuery = query.Encode()

	requestErr := Request(req, &trendings)

	if requestErr.IsError() {
		return trendings, requestErr
	}

	// save videoDetails to cache
	str, err := json.Marshal(trendings)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trendings, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	go redis_client.Client.SetKey(trendingCachePrefix, string(str), trendingCacheTTL)
	return trendings, errx.APIError{}
}

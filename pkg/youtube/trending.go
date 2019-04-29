package youtube

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// GetTrendings collects trendings from YouTube API
func GetTrendings() YoutubeResponse {
	trendings := YoutubeResponse{}
	req, err := http.NewRequest("GET", YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("chart", "mostPopular")
	req.URL.RawQuery = query.Encode()

	Request(req, &trendings)

	return trendings
}

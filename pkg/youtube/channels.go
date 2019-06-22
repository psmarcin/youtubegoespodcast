package youtube

import (
	"net/http"
	"ytg/pkg/errx"

	"github.com/sirupsen/logrus"
)

// GetChannels collects trendings from YouTube API
func GetChannels(q string) (YoutubeResponse, errx.APIError) {
	trendings := YoutubeResponse{}
	req, err := http.NewRequest("GET", YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return trendings, errx.NewAPIError(err, http.StatusBadRequest)
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("type", "channel")
	query.Add("q", q)
	query.Add("fields", "items(snippet(channelId,channelTitle,thumbnails/high,title))")
	query.Add("order", "viewCount")
	query.Add("maxResults", "5")
	req.URL.RawQuery = query.Encode()

	requestErr := Request(req, &trendings)
	if requestErr.IsError() {
		return trendings, requestErr
	}

	return trendings, errx.APIError{}
}

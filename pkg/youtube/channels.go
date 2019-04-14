package youtube

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetChannels collects trendings from YouTube API
func GetChannels(q string) YoutubeResponse {
	trendings := YoutubeResponse{}
	req, err := http.NewRequest("GET", YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("type", "channel")
	query.Add("q", q)
	query.Add("fields", "items(snippet(channelId,channelTitle,thumbnails/high,title))")
	query.Add("order", "viewCount")
	query.Add("maxResults", "5")
	req.URL.RawQuery = query.Encode()

	Request(req, &trendings)

	return trendings
}

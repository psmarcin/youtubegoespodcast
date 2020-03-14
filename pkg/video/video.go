package video

import (
	"strings"

	"github.com/iawia002/annie/extractors/youtube"
	"github.com/sirupsen/logrus"
)


// GetURL returns video URL based on videoID
func GetURL(videoID, format string) string {
	var foundUrl = ""
	var randomStreamKey = ""
	var url = youtubeBaseURL + "?v=" + videoID

	// extract details about video
	det, err := youtube.Extract(url)
	if err != nil {
		logrus.WithError(err).Info("[VIDEO]")
	}

	// try to find audio stream
	for _, detail := range det {
		for streamKey, stream := range detail.Streams{
			randomStreamKey = streamKey

			if strings.Contains(stream.Quality, format) {
				foundUrl = stream.URLs[0].URL
				return foundUrl
			}
		}
	}

	// fallback if no audio stream
	if foundUrl == "" && randomStreamKey != ""  {
		foundUrl = det[0].Streams[randomStreamKey].URLs[0].URL
	}

	return foundUrl
}

package video

import (
	"strings"

	"github.com/iawia002/annie/extractors/youtube"
	"github.com/sirupsen/logrus"
)

var audioFormats = []Format{
	FORMATS[141],
	FORMATS[140],
	FORMATS[139],
	FORMATS[171],
	FORMATS[251],
	FORMATS[250],
	FORMATS[249],
	FORMATS[18],
}

var audioFormat = "audio/mp4"
// GetURL returns video URL based on videoID
func GetURL(videoID string) string {
	var toReturn = ""
	var url = youtubeBaseURL + "?v=" + videoID
	det, err := youtube.Extract(url)
	if err != nil {
		logrus.WithError(err).Info("[VIDEO]")
	}

	for _, detail := range det {
		for _, stream := range detail.Streams{
			if(strings.Contains(stream.Quality, audioFormat)){
				toReturn = stream.URLs[0].URL
				return toReturn
			}
		}
	}

	return toReturn

	//for _, audioFormat := range audioFormats {
	//	for _, format := range details.Formats {
	//		url, err := getDownloadURL(format, details.htmlPlayerFile)
	//		if err != nil {
	//			logrus.WithError(err).Info("[VIDEO]")
	//			break
	//		}
	//		if url != nil && audioFormat.Itag == format.Itag {
	//			return url.String()
	//		}
	//	}
	//}
	// fallbackFormat := details.Formats[0]
	// url, err := getDownloadURL(fallbackFormat, details.htmlPlayerFile)
	// return url.String()
}

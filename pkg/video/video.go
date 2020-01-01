package video

import (
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

// GetURL returns video URL based on videoID
func GetURL(videoID string) string {
	var url = youtubeBaseURL + "?v=" + videoID
	det, err := youtube.Extract(url)
	// details, err := GetVideoInfoFromID(videoID)
	if err != nil {
		logrus.WithError(err).Info("[VIDEO]")
	}
	logrus.Infof("details %+v", det)

	// if len(detFormats) == 0 {
	// 	return ""
	// }
	return ""

	// for _, audioFormat := range audioFormats {
	// 	for _, format := range details.Formats {
	// 		url, err := getDownloadURL(format, details.htmlPlayerFile)
	// 		if err != nil {
	// 			logrus.WithError(err).Info("[VIDEO]")
	// 			break
	// 		}
	// 		if url != nil && audioFormat.Itag == format.Itag {
	// 			return url.String()
	// 		}
	// 	}
	// }
	// fallbackFormat := details.Formats[0]
	// url, err := getDownloadURL(fallbackFormat, details.htmlPlayerFile)
	// return url.String()
}

package video

import (
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
	details, err := GetVideoInfoFromID(videoID)
	if err != nil {
		logrus.WithError(err).Info("[VIDEO]")
	}

	for _, format := range details.Formats {
		for _, audioFormat := range audioFormats {
			url, err := getDownloadURL(format, details.htmlPlayerFile)
			if err != nil {
				logrus.WithError(err).Info("[VIDEO]")
				break
			}
			if url != nil && audioFormat.Itag == format.Itag {
				return url.String()
			}
		}
	}
	return ""
}

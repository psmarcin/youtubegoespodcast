package app

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/kkdai/youtube/v2"
	"go.opentelemetry.io/otel/attribute"
)

const (
	YoutubeVideoBaseURL   = "https://www.youtube.com/watch?v="
	DialTimeout           = 30
	DialKeepAlive         = 30
	IdleConnTimeout       = 60
	TLSHandshakeTimeout   = 10
	ExpectContinueTimeout = 1
)

type FileService struct{}

type Details struct {
	URL     url.URL
	Quality string
}

func NewFileService() FileService {
	return FileService{}
}

func (f FileService) GetDetails(ctx context.Context, videoID string) (Details, error) {
	details := Details{}
	_, span := tracer.Start(ctx, "file-get-details")
	span.SetAttributes(attribute.String("videoId", videoID))
	defer span.End()

	c := getYTClient()
	fullURL := YoutubeVideoBaseURL + videoID
	v, err := c.GetVideo(fullURL)
	if err != nil {
		return details, errors.Wrapf(err, "can't decode url: %s", fullURL)
	}

	audioStream, err := getAudioStream(v.Formats)
	if err != nil {
		return details, errors.Wrapf(err, "can't get audio stream url: %s", fullURL)
	}
	audioStreamRawURL := audioStream.URL

	audioStreamURL, err := url.Parse(audioStreamRawURL)
	if err != nil {
		return details, errors.Wrapf(err, "can't parse audio stream url: %s", audioStreamRawURL)
	}
	details.URL = *audioStreamURL

	return details, nil
}

func getAudioStream(formats []youtube.Format) (youtube.Format, error) {
	for _, format := range formats {
		if strings.HasPrefix(format.MimeType, "audio") {
			return format, nil
		}
	}

	if len(formats) > 0 {
		return formats[0], nil
	}

	return youtube.Format{}, errors.New("can't find audio")
}

func getHTTPClient() *http.Client {
	// setup a http client
	httpTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   DialTimeout * time.Second,
			KeepAlive: DialKeepAlive * time.Second,
		}).DialContext,
		IdleConnTimeout:       IdleConnTimeout * time.Second,
		TLSHandshakeTimeout:   TLSHandshakeTimeout * time.Second,
		ExpectContinueTimeout: ExpectContinueTimeout * time.Second,
	}

	httpClient := &http.Client{Transport: httpTransport}

	return httpClient
}

func getYTClient() youtube.Client {
	httpClient := getHTTPClient()

	client := youtube.Client{
		Debug:      false,
		HTTPClient: httpClient,
	}

	return client
}

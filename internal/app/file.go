package app

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/kkdai/youtube"
	"github.com/kkdai/youtube/pkg/decipher"
	"go.opentelemetry.io/otel/attribute"
)

const (
	YoutubeVideoBaseURL = "https://www.youtube.com/watch?v="
	HTTPClientTimeout   = 3 * time.Second
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
	span.SetAttributes(attribute.String("videoID", videoID))
	defer span.End()

	yt := youtube.NewYoutube(true, true)
	fullURL := YoutubeVideoBaseURL + videoID
	err := yt.DecodeURL(fullURL)
	if err != nil {
		return details, errors.Wrapf(err, "can't decode url: %s", fullURL)
	}

	audioStream, err := getAudioStream(yt.GetStreamInfo().Streams)
	if err != nil {
		return details, errors.Wrapf(err, "can't get audio stream url: %s", fullURL)
	}

	audioStreamRawURL, err := getStreamURL(videoID, audioStream)
	if err != nil {
		return details, errors.Wrapf(err, "can't get stream url for stream: %+v, url: %s", audioStream, fullURL)
	}

	audioStreamURL, err := url.Parse(audioStreamRawURL)
	if err != nil {
		return details, errors.Wrapf(err, "can't parse audio stream url: %s", audioStreamRawURL)
	}
	details.URL = *audioStreamURL

	return details, nil
}

func getAudioStream(streams []youtube.Stream) (youtube.Stream, error) {
	for _, stream := range streams {
		if strings.HasPrefix(stream.MimeType, "audio") {
			return stream, nil
		}
	}

	if len(streams) > 0 {
		return streams[0], nil
	}

	return youtube.Stream{}, errors.New("can't find audio")
}

func getStreamURL(videoID string, stream youtube.Stream) (string, error) {
	streamURL := stream.URL
	if streamURL == "" {
		cipher := stream.Cipher
		if cipher == "" {
			return "", errors.New("no cipher")
		}
		client := getHTTPClient()
		d := decipher.NewDecipher(client)
		decipherURL, err := d.Url(videoID, cipher)
		if err != nil {
			return "", err
		}
		streamURL = decipherURL
	}
	return streamURL, nil
}

func getHTTPClient() *http.Client {
	// setup a http client
	httpClient := &http.Client{
		Timeout: HTTPClientTimeout,
	}

	return httpClient
}

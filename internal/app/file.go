package app

import (
	"context"
	"net"
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

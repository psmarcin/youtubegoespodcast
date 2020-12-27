package app

import (
	"context"
	"errors"
	"github.com/kkdai/youtube"
	"go.opentelemetry.io/otel/label"
	"net/url"
	"strings"
)

const (
	YoutubeVideoBaseURL = "https://www.youtube.com/watch?v="
)

type FileService struct{}

type Details struct {
	Url     url.URL
	Quality string
}

func NewFileService() FileService {
	return FileService{}
}

func (f FileService) GetDetails(ctx context.Context, videoId string) (Details, error) {
	details := Details{}
	ctx, span := tracer.Start(ctx, "file-get-details")
	span.SetAttributes(label.String("videoId", videoId))
	defer span.End()

	yt := youtube.NewYoutube(true, true)
	err := yt.DecodeURL(YoutubeVideoBaseURL + videoId)
	if err != nil {
		return details, err
	}

	audioStream, err := getAudioStream(yt.GetStreamInfo().Streams)
	if err != nil {
		return details, err
	}

	audioStreamUrl, err := url.Parse(audioStream.URL)
	if err != nil {
		return details, err
	}
	details.Url = *audioStreamUrl

	return details, nil
}

func getAudioStream(streams []youtube.Stream) (youtube.Stream, error) {
	for _, stream := range streams {
		if strings.HasPrefix(stream.MimeType, "audio" ) {
			return stream, nil
		}
	}

	if len(streams) > 0 {
		return streams[0], nil
	}

	return youtube.Stream{}, errors.New("can't find audio")
}

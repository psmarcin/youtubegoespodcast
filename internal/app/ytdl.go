package app

import (
	"context"
	"github.com/pkg/errors"
	"github.com/psmarcin/youtubegoespodcast/internal/domain/video"
	"github.com/rylio/ytdl"
	"go.opentelemetry.io/otel/label"
	"net/url"
	"time"
)

const (
	FormatsNotFound     = "no formats available that match criteria"
	YoutubeVideoBaseURL = "https://www.youtube.com/watch?v="
)

type ytdlRepository interface {
	GetDownloadURL(context.Context, *ytdl.VideoInfo, *ytdl.Format) (*url.URL, error)
	GetVideoInfo(cx context.Context, value interface{}) (*ytdl.VideoInfo, error)
}

type YTDLService struct {
	ytdlRepository ytdlRepository
}

type YTDLVideo struct {
	ID             string
	URL            *url.URL
	FileUrl        *url.URL
	Description    string
	Title          string
	DatePublished  time.Time
	Keywords       []string
	Author         string
	Duration       time.Duration
	ContentType    string
	ContentLength  int64
	rawInformation *ytdl.VideoInfo
}

func NewYTDLService(ytdl ytdlRepository) YTDLService {
	return YTDLService{
		ytdlRepository: ytdl,
	}
}

func (y YTDLService) GetFileURL(ctx context.Context, info *ytdl.VideoInfo) (url.URL, error) {
	ctx, span := tracer.Start(ctx, "ytdl-get-file-url")
	span.SetAttributes(label.String("info", info.ID))
	defer span.End()

	var u url.URL

	formats := info.Formats
	filters := []string{
		"audio-only",
	}

	// parse filter arguments, and filter through formats
	for _, filter := range filters {
		filter, err := video.ParseFilter(filter)
		if err == nil {
			formats = filter(formats)
		}
	}

	if len(formats) == 0 {
		err := errors.New(FormatsNotFound)
		span.RecordError(ctx, err)
		return u, err
	}

	bestFormat := formats[0]
	fileURL, err := y.ytdlRepository.GetDownloadURL(context.Background(), info, bestFormat)
	if err != nil {
		err := errors.Wrapf(err, "can't get download url")
		span.RecordError(ctx, err)
		return u, err
	}

	return *fileURL, err
}

func (y YTDLService) GetFileInformation(ctx context.Context, videoId string) (YTDLVideo, error) {
	ctx, span := tracer.Start(ctx, "ytdl-get-file-information")
	span.SetAttributes(label.String("videoId", videoId))
	defer span.End()

	video := YTDLVideo{}

	u, err := url.Parse(YoutubeVideoBaseURL + videoId)
	if err != nil {
		return video, err
	}

	video.URL = u

	info, err := y.ytdlRepository.GetVideoInfo(context.Background(), videoId)
	if err != nil {
		span.RecordError(ctx, err)
		return video, err
	}
	video.ID = videoId
	video.rawInformation = info
	video.Title = info.Title
	video.Description = info.Description
	video.DatePublished = info.DatePublished
	video.Keywords = info.Keywords
	video.Author = info.Uploader
	video.Duration = info.Duration

	fileUrl, err := y.GetFileURL(ctx, info)
	if err != nil {
		span.RecordError(ctx, err)
		return video, err
	}
	video.FileUrl = &fileUrl

	//fileDetails, err := GetFileDetails(video.FileUrl, videoID)
	//if err != nil {
	//	l.WithError(err).Errorf("can't get file video for video %s", videoID)
	//}
	//
	//video.ContentLength = fileDetails.ContentLength
	//video.ContentType = fileDetails.ContentType
	return video, nil
}

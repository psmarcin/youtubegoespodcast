package video

import (
	"context"
	"errors"
	"fmt"
	"github.com/rylio/ytdl"
	"net/url"
	"strings"
	"time"
)

const (
	DetailsNotFound     = "rawInformation not provided, don't have from get file url"
	FormatsNotFound     = "no formats available that match criteria"
	YoutubeVideoBaseURL = "https://www.youtube.com/watch?v="
)

type Video struct {
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

type FileUrlGetter interface {
	GetDownloadURL(context.Context, *ytdl.VideoInfo, *ytdl.Format) (*url.URL, error)
}

type FileInformationGetter interface {
	GetVideoInfo(cx context.Context, value interface{}) (*ytdl.VideoInfo, error)
}

func New(ID string) Video {
	u, _ := url.Parse(YoutubeVideoBaseURL + ID)
	v := Video{
		ID:  ID,
		URL: u,
	}
	return v
}

func (v Video) GetFileURL(fileUrlGetter FileUrlGetter) (url.URL, error) {
	var u url.URL

	if v.rawInformation == nil {
		return u, errors.New(DetailsNotFound)
	}

	formats := v.rawInformation.Formats
	filters := []string{
		"audio-only",
	}

	// parse filter arguments, and filter through formats
	for _, filter := range filters {
		filter, err := parseFilter(filter)
		if err == nil {
			formats = filter(formats)
		}
	}

	if len(formats) == 0 {
		return u, errors.New(FormatsNotFound)
	}

	bestFormat := formats[0]
	fileURL, err := fileUrlGetter.GetDownloadURL(context.Background(), v.rawInformation, bestFormat)

	return *fileURL, err
}

func (v Video) GetFileInformation(fileInformationGetter FileInformationGetter, fileUrlGetter FileUrlGetter) (Video, error) {
	video := v

	info, err := fileInformationGetter.GetVideoInfo(context.Background(), video.URL)
	if err != nil {
		return video, err
	}
	video.rawInformation = info
	video.Title = info.Title
	video.Description = info.Description
	video.DatePublished = info.DatePublished
	video.Keywords = info.Keywords
	video.Author = info.Uploader
	video.Duration = info.Duration

	fileUrl, err := video.GetFileURL(fileUrlGetter)
	if err != nil {
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

func parseFilter(filterString string) (func(ytdl.FormatList) ytdl.FormatList, error) {
	filterString = strings.TrimSpace(filterString)

	switch filterString {
	case "audio-only":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.
				Extremes(ytdl.FormatResolutionKey, filterString == "")
		}, nil
	}

	err := fmt.Errorf("Invalid filter")
	split := strings.SplitN(filterString, ":", 2)
	if len(split) != 2 {
		return nil, err
	}

	key := ytdl.FormatKey(split[0])
	exclude := key[0] == '!'
	if exclude {
		key = key[1:]
	}

	value := strings.TrimSpace(split[1])
	if value == "best" || value == "worst" {
		return func(formats ytdl.FormatList) ytdl.FormatList {
			f := formats.Extremes(key, value == "best")
			if exclude {
				f = formats.Subtract(f)
			}
			return f
		}, nil
	}

	vals := strings.Split(value, ",")
	values := make([]interface{}, len(vals))

	for i, v := range vals {
		values[i] = strings.TrimSpace(v)
	}

	return func(formats ytdl.FormatList) ytdl.FormatList {
		f := formats.Filter(key, values)
		if exclude {
			f = formats.Subtract(f)
		}
		return f
	}, nil
}

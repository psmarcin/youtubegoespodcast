package video

import (
	"context"
	"fmt"
	"github.com/psmarcin/ytdl"
	"net/url"
	"strings"
	"time"
)

var client = ytdl.DefaultClient

type Details struct {
	FileUrl       *url.URL
	Description   string
	Title         string
	DatePublished time.Time
	Keywords      []string
	Author        string
	Duration      time.Duration
}

// GetDetails returns video details based on videoID
func GetDetails(videoID string) (Details, error) {
	var details Details

	info, err := client.GetVideoInfo(context.Background(), videoID)
	if err != nil {
		err = fmt.Errorf("Unable to fetch video info: %s", err)
		return details, err
	}

	details.Title = info.Title
	details.Description = info.Description
	details.DatePublished = info.DatePublished
	details.Keywords = info.Keywords
	details.Author = info.Uploader
	details.Duration = info.Duration

	details.FileUrl, err = getVideoUrl(info)
	return details, err
}

func getVideoUrl(info *ytdl.VideoInfo) (*url.URL, error) {
	var u *url.URL
	formats := info.Formats
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
		err := fmt.Errorf("No formats available that match criteria")
		return u, err
	}

	bestFormat := formats[0]
	u, err := client.GetDownloadURL(context.Background(), info, bestFormat)

	return u, err
}

func parseFilter(filterString string) (func(ytdl.FormatList) ytdl.FormatList, error) {

	filterString = strings.TrimSpace(filterString)
	switch filterString {
	case "best", "worst":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatResolutionKey, filterString == "best").Extremes(ytdl.FormatAudioBitrateKey, filterString == "best")
		}, nil
	case "best-fps", "worst-fps":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatFPSKey, filterString == "best-fps").Extremes(ytdl.FormatResolutionKey, filterString == "best-fps").Extremes(ytdl.FormatAudioBitrateKey, filterString == "best-fps")
		}, nil
	case "best-video", "worst-video":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatResolutionKey, strings.HasPrefix(filterString, "best"))
		}, nil
	case "best-audio", "worst-audio":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.Extremes(ytdl.FormatAudioBitrateKey, strings.HasPrefix(filterString, "best"))
		}, nil
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

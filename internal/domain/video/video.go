package video

import (
	"github.com/pkg/errors"
	"github.com/rylio/ytdl"
	"strings"
)

func ParseFilter(filterString string) (func(ytdl.FormatList) ytdl.FormatList, error) {
	filterString = strings.TrimSpace(filterString)

	switch filterString {
	case "audio-only":
		return func(formats ytdl.FormatList) ytdl.FormatList {
			return formats.
				Extremes(ytdl.FormatResolutionKey, filterString == "")
		}, nil
	}

	err := errors.New("invalid filter")
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

package video

import (
	"strconv"
	"strings"
	"testing"
)

func TestGetURL(t *testing.T) {
	type args struct {
		videoID string
	}
	type want struct {
		url  string
		itag string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "videoID: tGL4tKXkA4U",
			args: args{
				videoID: "tGL4tKXkA4U",
			},
			want: want{
				url:  "googlevideo.com/videoplayback",
				itag: "140",
			},
		},
		{
			name: "videoID: oIbEvwzLxgE",
			args: args{
				videoID: "oIbEvwzLxgE",
			},
			want: want{
				url:  "googlevideo.com/videoplayback",
				itag: "140",
			},
		},
		{
			name: "videoID: oIbEvwzLxgE1",
			args: args{
				videoID: "oIbEvwzLxgE1",
			},
			want: want{
				url:  "googlevideo.com/videoplayback",
				itag: "140",
			},
		},
		{
			name: "videoID: uwfdFCP3KYM",
			args: args{
				videoID: "uwfdFCP3KYM",
			},
			want: want{
				url:  "googlevideo.com/videoplayback",
				itag: "140",
			},
		},
		{
			name: "videoID: jNQXAC9IVRw",
			args: args{
				videoID: "jNQXAC9IVRw",
			},
			want: want{
				url:  "googlevideo.com/videoplayback",
				itag: "140",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetURL(tt.args.videoID)
			containsURL := strings.Contains(got, tt.want.url)
			containsItag := strings.Contains(got, "itag="+tt.want.itag)
			foundFormat := false
			for _, audioFormat := range audioFormats {
				if strings.Contains(got, strconv.Itoa(audioFormat.Itag)) {
					foundFormat = true
					break
				}
			}
			if !containsURL || !foundFormat || !containsItag {
				t.Errorf("GetURL() = %v, want url: %v, formatFound %t", got, tt.want.url, foundFormat)
			}
		})
	}
}

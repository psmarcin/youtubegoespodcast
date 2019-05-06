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
		url string
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
				url: "googlevideo.com/videoplayback",
			},
		},
		{
			name: "videoID: oIbEvwzLxgE",
			args: args{
				videoID: "oIbEvwzLxgE",
			},
			want: want{
				url: "googlevideo.com/videoplayback",
			},
		},
		{
			name: "videoID: oIbEvwzLxgE1",
			args: args{
				videoID: "oIbEvwzLxgE1",
			},
			want: want{
				url: "googlevideo.com/videoplayback",
			},
		},
		{
			name: "videoID: uwfdFCP3KYM",
			args: args{
				videoID: "uwfdFCP3KYM",
			},
			want: want{
				url: "googlevideo.com/videoplayback",
			},
		},
		{
			name: "videoID: _kMShzhW3KE",
			args: args{
				videoID: "_kMShzhW3KE",
			},
			want: want{
				url: "googlevideo.com/videoplayback",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetURL(tt.args.videoID)
			containsURL := strings.Contains(got, tt.want.url)
			foundFormat := false
			for _, audioFormat := range audioFormats {
				if strings.Contains(got, strconv.Itoa(audioFormat.Itag)) {
					foundFormat = true
					break
				}
			}
			if !containsURL || !foundFormat {
				t.Errorf("GetURL() = %v, want url: %v, formatFound %t", got, tt.want.url, foundFormat)
			}
		})
	}
}

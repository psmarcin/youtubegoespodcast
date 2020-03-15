package video

import (
	"testing"
)

func TestGetURL(t *testing.T) {
	type args struct {
		videoID string
		format  string
	}
	type want struct {
		url  string
		itag string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "videoID: tGL4tKXkA4U",
			args: args{
				videoID: "tGL4tKXkA4U",
				format:  FormatMp4,
			},
		},
		{
			name: "videoID: uwfdFCP3KYM",
			args: args{
				videoID: "uwfdFCP3KYM",
				format:  FormatMp4,
			},
		},
		{
			name: "videoID: jNQXAC9IVRw",
			args: args{
				videoID: "jNQXAC9IVRw",
				format:  FormatMp4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetURL(tt.args.videoID, tt.args.format)

			if got == "" {
				t.Errorf("GetURL() = %v, want url not empty", got)
			}
		})
	}
}

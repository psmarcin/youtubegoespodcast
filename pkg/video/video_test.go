package video

import (
	"testing"
)

func TestGetURL(t *testing.T) {
	type args struct {
		videoID string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "videoID: tGL4tKXkA4U",
			args: args{
				videoID: "tGL4tKXkA4U",
			},
		},
		{
			name: "videoID: uwfdFCP3KYM",
			args: args{
				videoID: "uwfdFCP3KYM",
			},
		},
		{
			name: "videoID: jNQXAC9IVRw",
			args: args{
				videoID: "jNQXAC9IVRw",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetURL(tt.args.videoID)

			if got == "" || err != nil {
				t.Errorf("GetURL() = %v, want url not empty", got)
			}
		})
	}
}

func TestGetURLGetNotExistedVideo(t *testing.T) {
	notExistingVideoID := "xxx"
	got, err := GetURL(notExistingVideoID)

	if got != "" && err == nil {
		t.Errorf("GetURL() = %v, want url to be empty", got)
	}
}

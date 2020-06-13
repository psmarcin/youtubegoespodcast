package video

import (
	"github.com/rylio/ytdl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var deps = Dependencies{
	Details:    HeadRequest,
	Info:       ytdl.GetVideoInfo,
	GetFileUrl: GetVideoUrl,
}

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
			got, err := GetDetails(tt.args.videoID, deps)

			if got.FileUrl.String() == "" || err != nil {
				t.Errorf("GetDetails() = %v, want url not empty", got)
			}
		})
	}
}

func TestGetDetails(t *testing.T) {
	type args struct {
		videoID string
	}
	tests := []struct {
		name    string
		args    args
		want    Details
		wantErr bool
	}{
		{
			name: "should get duration",
			args: args{
				videoID: "jNQXAC9IVRw",
			},
			want: Details{
				Title:    "Me at the zoo",
				Author:   "jawed",
				Duration: 19 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "should get error on videoId not found",
			args: args{
				videoID: "xxx",
			},
			want:    Details{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDetails(tt.args.videoID, deps)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Duration, got.Duration)
			assert.Equal(t, tt.want.Title, got.Title)
			assert.Equal(t, tt.want.Author, got.Author)
		})
	}
}

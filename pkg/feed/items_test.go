package feed

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func Test_getVideoFileDetails(t *testing.T) {
	type args struct {
		videoURL string
	}
	tests := []struct {
		name    string
		args    args
		want    VideoFileDetails
		wantErr bool
	}{
		{
			name: "get video file details - found",
			args: args{
				videoURL: os.Getenv("API_URL") + "video/123.mp3",
			},
			want: VideoFileDetails{
				ContentType:   "app",
				ContentLength: "123",
			},
			wantErr: false,
		},
		{
			name: "get video file details - not found",
			args: args{
				videoURL: os.Getenv("API_URL") + "video/123.mp3",
			},
			want:    VideoFileDetails{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution
			status := http.StatusOK
			if tt.wantErr == true {
				status = http.StatusNotFound
			}
			t.Logf("Status %d", status)
			gock.New(os.Getenv("API_URL")).
				Head("video/123.mp3").
				Reply(status).
				AddHeader("Content-Type", tt.want.ContentType).
				AddHeader("Content-Length", tt.want.ContentLength)

			got, err := GetVideoFileDetails(tt.args.videoURL)
			if (err.IsError()) != tt.wantErr {
				t.Errorf("GetVideoFileDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVideoFileDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

package feed

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"gopkg.in/h2non/gock.v1"
)

func Test_parseDuration(t *testing.T) {
	want1, _ := time.ParseDuration("1h0m0s")

	type args struct {
		duration string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "parse 1h0m0s",
			args: args{
				duration: "1h0m0s",
			},
			want:    want1,
			wantErr: false,
		},
		{
			name: "parse PT1h0m0s",
			args: args{
				duration: "PT1h0m0s",
			},
			want:    want1,
			wantErr: false,
		},
		{
			name: "parse PT1H0M0S",
			args: args{
				duration: "PT1H0M0S",
			},
			want:    want1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getVideoDetails(t *testing.T) {
	duration, _ := time.ParseDuration("1h0m0s")

	type args struct {
		videoID string
	}
	tests := []struct {
		name    string
		args    args
		want    VideoDetails
		wantErr bool
	}{
		{
			name: "get video duration",
			args: args{
				videoID: "123",
			},
			want: VideoDetails{
				Duration: duration,
			},
			wantErr: false,
		},
		{
			name: "get error on bad request response from google api",
			args: args{
				videoID: "ERR",
			},
			want:    VideoDetails{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution
			status := http.StatusOK
			response := `{
					"items": [
						{
							"contentDetails": {
								"duration": "1h0m0s"
							}
						}
					]
				}
				`
			if tt.args.videoID == "ERR" {
				status = http.StatusBadRequest
				response = ``
			}
			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/videos").
				Reply(status).
				BodyString(response)

			got, err := GetVideoDetails(tt.args.videoID, true)
			if (err.IsError()) != tt.wantErr {
				t.Errorf("GetVideoDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVideoDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

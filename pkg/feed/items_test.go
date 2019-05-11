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

			got, err := getVideoDetails(tt.args.videoID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVideoDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVideoDetails() = %v, want %v", got, tt.want)
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

			got, err := getVideoFileDetails(tt.args.videoURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVideoFileDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVideoFileDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeed_getVideos(t *testing.T) {
	type fields struct {
		XMLName       string
		ChannelID     string
		Title         string
		Link          string
		Description   string
		Category      string
		Generator     string
		Language      string
		LastBuildDate string
		PubDate       string
		Image         Image
		ITAuthor      string
		ITSubtitle    string
		ITSummary     ITSummary
		ITImage       ITImage
		ITExplicit    string
		ITCategory    ITCategory
		Items         []Item
	}
	tests := []struct {
		name    string
		fields  fields
		want    VideosResponse
		wantErr bool
	}{
		{
			name: "get videos",
			fields: fields{
				ChannelID: "123",
			},
			want: VideosResponse{
				Items: []VideosItems{
					{
						ID: ID{
							VideoID: "123",
						},
						Snippet: VideosSnippet{
							ChannelID: "channel",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{
				XMLName:       tt.fields.XMLName,
				ChannelID:     tt.fields.ChannelID,
				Title:         tt.fields.Title,
				Link:          tt.fields.Link,
				Description:   tt.fields.Description,
				Category:      tt.fields.Category,
				Generator:     tt.fields.Generator,
				Language:      tt.fields.Language,
				LastBuildDate: tt.fields.LastBuildDate,
				PubDate:       tt.fields.PubDate,
				Image:         tt.fields.Image,
				ITAuthor:      tt.fields.ITAuthor,
				ITSubtitle:    tt.fields.ITSubtitle,
				ITSummary:     tt.fields.ITSummary,
				ITImage:       tt.fields.ITImage,
				ITExplicit:    tt.fields.ITExplicit,
				ITCategory:    tt.fields.ITCategory,
				Items:         tt.fields.Items,
			}
			defer gock.Off() // Flush pending mocks after test execution
			status := http.StatusOK
			response := `{
					"items": [
						{
							"id": {
								"videoId": "123"
							},
							"snippet": {
								"channelId": "channel"
							}
						}
					]
				}
				`
			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/search").
				Reply(status).
				BodyString(response)
			got, err := f.getVideos("123")
			if (err != nil) != tt.wantErr {
				t.Errorf("Feed.getVideos() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Feed.getVideos() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_calculateDuration(t *testing.T) {
	type args struct {
		d time.Duration
	}

	duration1h00m10s, _ := time.ParseDuration("1h00m10s")
	duration123h59m59s, _ := time.ParseDuration("123h59m59s")

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1h00m10s",
			args: args{
				d: duration1h00m10s,
			},
			want: "01:00:10",
		},
		{
			name: "123h59m59s",
			args: args{
				d: duration123h59m59s,
			},
			want: "123:59:59",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDuration(tt.args.d); got != tt.want {
				t.Errorf("calculateDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

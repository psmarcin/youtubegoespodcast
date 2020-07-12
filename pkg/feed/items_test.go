package feed

import (
	"context"
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/video"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/rylio/ytdl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"
)

type MockVideoDependencies struct {
	mock.Mock
}

func (m *MockVideoDependencies) Details(videoID string) (*http.Response, error) {
	args := m.Called(videoID)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockVideoDependencies) Info(cx context.Context, value interface{}) (*ytdl.VideoInfo, error) {
	args := m.Called(cx, value)
	return args.Get(0).(*ytdl.VideoInfo), args.Error(1)
}

func TestFeed_SetVideos(t *testing.T) {
	i := Item{
		Video: youtube.Video{
			ID:          "ivid1",
			Title:       "ivt1",
			Description: "ivd1",
			PublishedAt: time.Now(),
			Url:         "ivu1",
		},
		Details: video.Video{
			Description:   "idd1",
			Title:         "idt1",
			DatePublished: time.Now().Truncate(time.Hour * 24),
			Author:        "ida1",
			Duration:      100,
			ContentType:   "idct1",
			ContentLength: 1000,
		},
	}

	f := &Feed{
		ChannelID: "123",
		Content:   podcast.Podcast{},
	}
	f.Items = []Item{
		i,
	}
	err := f.SetVideos()

	assert.NoError(t, err)

	assert.Len(t, f.Content.Items, 1)
	assert.Equal(t, f.Content.Items[0].Title, i.Details.Title)
	assert.Equal(t, f.Content.Items[0].GUID, i.Video.Url)
	assert.Equal(t, f.Content.Items[0].Link, i.Video.Url)
	assert.Equal(t, f.Content.Items[0].Description, i.Details.Description)
	assert.Equal(t, f.Content.Items[0].PubDate, &i.Details.DatePublished)
	assert.Contains(t, f.Content.Items[0].Enclosure.URL, i.Video.ID+"/track.mp3")
	assert.Equal(t, f.Content.Items[0].Enclosure.Type, podcast.MP3)
	assert.Equal(t, f.Content.LastBuildDate, i.Details.DatePublished.Format(time.RFC1123Z))
}

func Test_filterEmptyVideoOut(t *testing.T) {
	type args struct {
		items []Item
	}
	tests := []struct {
		name string
		args args
		want []Item
	}{
		{
			name: "should remove all items",
			args: args{
				items: []Item{
					{},
					{},
					{},
					{},
				},
			},
			want: []Item{},
		},
		{
			name: "should remove all but one item",
			args: args{
				items: []Item{
					{
						Video: youtube.Video{
							ID: "123",
						},
					},
					{},
					{},
					{},
				},
			},
			want: []Item{
				{
					Video: youtube.Video{
						ID: "123",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterEmptyVideoOut(tt.args.items)
			assert.Len(t, got, len(tt.want))
			assert.Subset(t, got, tt.want)
		})
	}
}

func TestFeed_SetItems(t *testing.T) {
	type fields struct {
		ChannelID string
		Content   podcast.Podcast
		Items     []Item
	}
	type args struct {
		videos []youtube.Video
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "should set all items",
			args: args{
				videos: []youtube.Video{
					{
						ID: "123",
					},
					{},
					{},
					{},
				},
			},
		},
		{
			name: "should set no items",
			args: args{
				videos: []youtube.Video{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{
				ChannelID: tt.fields.ChannelID,
				Content:   tt.fields.Content,
				Items:     tt.fields.Items,
			}
			f.SetItems(tt.args.videos)
			assert.Len(t, f.Items, len(tt.args.videos))
		})
	}
}

//func TestFeed_EnrichItems(t *testing.T) {
//	deps := new(MockVideoDependencies)
//	resp := httptest.NewRecorder()
//	resp.Header().Set("Content-Type", "mp3")
//	resp.Header().Set("Content-Length", "444444")
//
//	deps.On("Details",  "i1").Return(resp.Result(), nil)
//	deps.On("Info",  context.Background(), "i1").Return(&ytdl.VideoInfo{
//		ID:          "i1",
//		Title:       "t1",
//		Description: "d1",
//		Formats: ytdl.FormatList{
//			&ytdl.Format{
//				Adaptive: true,
//				Itag:
//			},
//		},
//	}, nil)
//
//	d := video.Dependencies{
//		Details: deps.Details,
//		Info:    deps.Info,
//	}
//	f := &Feed{
//		ChannelID: "ch1",
//		Items: []Item{
//			Item{
//				Video:   youtube.Video{
//					ID: "i1",
//				},
//				Details: video.Details{},
//			},
//		},
//	}
//	err := f.EnrichItems(d)
//	assert.NoError(t, err)
//	assert.Equal(t, f.Items[0].Details.ContentType, "mp3")
//	assert.Equal(t, f.Items[0].Details.ContentLength, "444444")
//}

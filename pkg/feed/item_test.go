package feed

import (
	"errors"
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/video"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

const (
	YoutubeVideoBaseURL = "https://www.youtube.com/watch?v="
)

func TestFeed_SetVideos(t *testing.T) {
	id := "123"
	u, _ := url.Parse(YoutubeVideoBaseURL + id)

	i := Item{
		ID:          "ivid1",
		Title:       "ivt1",
		Description: "ivd1",
		PubDate:     time.Now(),
		Link:        u,
		GUID:        u.String(),
		FileLength:  int64(123),
		Duration:    time.Hour,
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
	assert.Equal(t, f.Content.Items[0].Title, i.Title)
	assert.Equal(t, f.Content.Items[0].GUID, i.Link.String())
	assert.Equal(t, f.Content.Items[0].Link, i.Link.String())
	assert.Equal(t, f.Content.Items[0].Description, i.Description)
	assert.Equal(t, f.Content.Items[0].PubDate, &i.PubDate)
	assert.Contains(t, f.Content.Items[0].Enclosure.URL, i.ID+"/track.mp3")
	assert.Equal(t, f.Content.Items[0].Enclosure.Type, podcast.MP3)
	assert.Equal(t, f.Content.LastBuildDate, i.PubDate.Format(time.RFC1123Z))
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
						ID:          "JZAunPKoHL0",
						Url:         YoutubeVideoBaseURL + "JZAunPKoHL0",
						Title:       "t1",
						Description: "d2",
					},
					{
						ID:          "JZAunPKoHL0",
						Url:         YoutubeVideoBaseURL + "JZAunPKoHL0",
						Title:       "t1",
						Description: "d2",
					},
					{
						ID:          "JZAunPKoHL0",
						Url:         YoutubeVideoBaseURL + "JZAunPKoHL0",
						Title:       "t1",
						Description: "d2",
					},
					{
						ID:          "JZAunPKoHL0",
						Url:         YoutubeVideoBaseURL + "JZAunPKoHL0",
						Title:       "t1",
						Description: "d2",
					},
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
			err := f.SetItems(tt.args.videos)
			assert.NoError(t, err)
			assert.Len(t, f.Items, len(tt.args.videos))
		})
	}
}

func TestItemV2_IsValid(t *testing.T) {
	type fields struct {
		Title       string
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should return false for no title",
			fields: fields{
				Title:       "",
				Description: "desc",
			},
			want: false,
		},
		{
			name: "should return false for no description",
			fields: fields{
				Title:       "title",
				Description: "",
			},
			want: false,
		},
		{
			name: "should return false for no description and no title",
			fields: fields{
				Title:       "",
				Description: "",
			},
			want: false,
		},
		{
			name: "should return true for description and title",
			fields: fields{
				Title:       "title",
				Description: "description",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := Item{
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
			}
			if got := item.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItem_enrichPopulateAllDate(t *testing.T) {
	id := "123"
	u, _ := url.Parse(YoutubeVideoBaseURL + id)
	fn := func(id string, fig video.FileInformationGetter, fug video.FileUrlGetter) (video.Video, error) {
		return video.Video{
			ID:            id,
			Title:         "title",
			Duration:      time.Hour,
			FileUrl:       u,
			ContentLength: 10000,
		}, nil
	}

	item := &Item{ID: "123"}
	expectedItem := &Item{
		ID:         "123",
		FileURL:    u,
		Duration:   time.Hour,
		FileLength: 10000,
	}
	err := item.enrich(fn)
	assert.NoError(t, err)
	assert.EqualValues(t, item, expectedItem)
}

func TestItem_enrichReturnError(t *testing.T) {
	expectedErr := errors.New("test error")
	fn := func(id string, fig video.FileInformationGetter, fug video.FileUrlGetter) (video.Video, error) {
		return video.Video{}, expectedErr
	}

	item := &Item{ID: "123"}
	err := item.enrich(fn)
	assert.Error(t, err, err)
}

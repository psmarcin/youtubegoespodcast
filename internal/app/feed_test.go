package app

import (
	"context"
	"github.com/eduncan911/podcast"
	"github.com/pkg/errors"
	"github.com/rylio/ytdl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type YouTubeDependencyMock struct {
	mock.Mock
}

func (m *YouTubeDependencyMock) ListEntry(s string) ([]YouTubeFeedEntry, error) {
	args := m.Called(s)
	return args.Get(0).([]YouTubeFeedEntry), args.Error(1)
}

func (m *YouTubeDependencyMock) GetChannelCache(s string) (YouTubeChannel, error) {
	args := m.Called(s)
	return args.Get(0).(YouTubeChannel), args.Error(1)
}

type YTDLDependencyMock struct {
	mock.Mock
}

func (m *YTDLDependencyMock) GetFileURL(info *ytdl.VideoInfo) (url.URL, error) {
	args := m.Called(info)
	return args.Get(0).(url.URL), args.Error(1)
}

func (m *YTDLDependencyMock) GetFileInformation(videoId string) (YTDLVideo, error) {
	args := m.Called(videoId)
	return args.Get(0).(YTDLVideo), args.Error(1)
}

type DependencyVideoMock struct {
	mock.Mock
}

func (d *DependencyVideoMock) Details(videoID string) (*http.Response, error) {
	args := d.Called(videoID)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (d *DependencyVideoMock) Info(cx context.Context, value interface{}) (*ytdl.VideoInfo, error) {
	args := d.Called(cx, value)
	return args.Get(0).(*ytdl.VideoInfo), args.Error(1)
}

func (d *DependencyVideoMock) GetFileUrl(info *ytdl.VideoInfo) (*url.URL, error) {
	args := d.Called(info)
	return args.Get(0).(*url.URL), args.Error(1)
}

type YT struct {
	returnChannel YouTubeChannel
	returnError   error
}

func TestFeed_serialize(t *testing.T) {
	ytM := new(YouTubeDependencyMock)
	ytM.On("ListEntry", "123").Return([]YouTubeFeedEntry{}, nil)
	ytM.On("GetChannelCache", "123").Return(YouTubeChannel{}, nil)

	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileURL", "123").Return([]YouTubeFeedEntry{}, nil)
	ytdlM.On("GetFileInformation", "123").Return(YouTubeChannel{}, nil)

	feedService := NewFeedService(ytM, ytdlM)
	f, _ := feedService.Create("123")
	ti := time.Now()
	f.Content = podcast.New("title", "http://onet", "description", &ti, &ti)
	serialized, err := f.Serialize()

	assert.NoError(t, err, "serialize should not return error")
	assert.Contains(t, string(serialized), "<link>http://onet</link>")
}

func TestCreateShouldReturnFeedWithVideo1(t *testing.T) {
	ytM := new(YouTubeDependencyMock)
	ytM.On("ListEntry", "123").Return([]YouTubeFeedEntry{{
		ID:          "vi1",
		Title:       "vt1",
		Description: "vd1",
	}}, nil)
	ytM.On("GetChannelCache", "123").Return(YouTubeChannel{
		ChannelId:   "123",
		Country:     "pl",
		Description: "d1",
		Title:       "t1",
		Url:         "u1",
	}, nil)

	vd := new(DependencyVideoMock)
	fileURLRaw := "http://onet.pl"
	fileURL, _ := url.Parse(fileURLRaw)
	info := YTDLVideo{
		ID:            "",
		Title:         "Title Video",
		Description:   "Description Video",
		DatePublished: time.Time{},
		Keywords:      nil,
		Duration:      0,
	}

	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileUrl", info).Return(fileURL, nil)
	ytdlM.On("GetFileInformation", "vi1").Return(info, nil)
	resp := httptest.NewRecorder()

	resp.Header().Set("Content-Type", "mp3")
	resp.Header().Set("Content-Length", "123")
	vd.On("Details", fileURLRaw).Return(resp.Result(), nil)

	feedService := NewFeedService(ytM, ytdlM)
	f, _ := feedService.Create("123")

	assert.Equal(t, "123", f.ChannelID)
	assert.Len(t, f.Content.Items, 1)
	assert.True(t, ytM.AssertNumberOfCalls(t, "ListEntry", 1))
	assert.True(t, ytM.AssertNumberOfCalls(t, "GetChannelCache", 1))
	assert.Equal(t, "pl", f.Content.Language)
	assert.Equal(t, "d1", f.Content.Description)
	assert.Equal(t, "t1", f.Content.Title)
	assert.Equal(t, "u1", f.Content.Link)
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
			item := FeedItem{
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
	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileInformation", id).Return(YTDLVideo{
		ID:            "JZAunPKoHL0",
		URL:           u,
		Description:   "d2",
		Title:         "title",
		Duration:      time.Hour,
		FileUrl:       u,
		ContentLength: 10000,
	}, nil)

	item := FeedItem{ID: id}
	expectedItem := FeedItem{
		ID:         id,
		FileURL:    u,
		Duration:   time.Hour,
		FileLength: 10000,
	}
	feedService := NewFeedService(YouTubeService{}, ytdlM)

	feedItem, err := feedService.Enrich(item)
	assert.NoError(t, err)
	assert.EqualValues(t, feedItem, expectedItem)
}

func TestItem_enrichReturnError(t *testing.T) {
	errorMessage := "error message"
	id := "123"
	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileInformation", id).Return(YTDLVideo{}, errors.New(errorMessage))

	item := FeedItem{ID: id}
	feedService := NewFeedService(YouTubeService{}, ytdlM)

	_, err := feedService.Enrich(item)
	assert.Error(t, err, err)
}

func (yt YT) GetChannelCache(_ string) (YouTubeChannel, error) {
	return yt.returnChannel, yt.returnError
}

func TestFeed_GetDetails(t *testing.T) {
	ytM := new(YouTubeDependencyMock)

	feedService := NewFeedService(ytM, YTDLService{})

	type fields struct {
		ChannelID string
		Content   podcast.Podcast
	}
	type args struct {
		channelID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		expect  podcast.Podcast
		wantErr bool
		before  func()
	}{
		{
			name: "should get channel details",
			fields: fields{
				ChannelID: "",
				Content:   podcast.Podcast{},
			},
			args: args{
				channelID: "1",
			},
			expect: podcast.Podcast{
				Title:       "t1",
				Link:        "u1",
				Description: "d1",
			},
			wantErr: false,
			before: func() {
				ytM.On("GetChannelCache", "1").Return(YouTubeChannel{
					ChannelId:   "ch1",
					Country:     "pl1",
					Description: "d1",
					PublishedAt: time.Date(2000, 11, 01, 1, 1, 1, 1, time.UTC),
					Thumbnail:   "th1",
					Title:       "t1",
					Url:         "u1",
				}, nil)
			},
		},
		{
			name: "should throw error on youtube service error",
			fields: fields{
				ChannelID: "",
				Content:   podcast.Podcast{},
			}, args: args{
			channelID: "2",
		},
			wantErr: true,
			before: func() {
				ytM.On("GetChannelCache", "2").Return(YouTubeChannel{}, errors.New("can't get channel"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{
				ChannelID: tt.fields.ChannelID,
				Content:   tt.fields.Content,
			}

			tt.before()

			pod, err := feedService.GetDetails(tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetails() error = %v, wantErr %v", err, tt.wantErr)
			}

			if pod.Title != tt.expect.Title &&
				pod.Description != tt.expect.Description &&
				pod.Link != tt.expect.Link {
				t.Errorf("GetDetails() get = %v, want = %v", f.Content, tt.expect)
			}
		})
	}
}

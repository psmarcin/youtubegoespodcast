package app

import (
	"context"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/eduncan911/podcast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type YouTubeDependencyMock struct {
	mock.Mock
}

func (m *YouTubeDependencyMock) ListEntry(ctx context.Context, s string) ([]YouTubeFeedEntry, error) {
	args := m.Called(mock.Anything, s)
	return args.Get(0).([]YouTubeFeedEntry), args.Error(1)
}

func (m *YouTubeDependencyMock) GetChannelCache(ctx context.Context, s string) (YouTubeChannel, error) {
	args := m.Called(mock.Anything, s)
	return args.Get(0).(YouTubeChannel), args.Error(1)
}

type YTDLDependencyMock struct {
	mock.Mock
}

func (m *YTDLDependencyMock) GetDetails(ctx context.Context, videoID string) (Details, error) {
	args := m.Called(mock.Anything, videoID)
	return args.Get(0).(Details), args.Error(1)
}

func TestFeed_serialize(t *testing.T) {
	ytM := new(YouTubeDependencyMock)
	ytM.On("ListEntry", context.Background(), "123").Return([]YouTubeFeedEntry{}, nil)
	ytM.On("GetChannelCache", context.Background(), "123").Return(YouTubeChannel{}, nil)

	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileURL", context.Background(), "123").Return([]YouTubeFeedEntry{}, nil)
	ytdlM.On("GetFileInformation", context.Background(), "123").Return(YouTubeChannel{}, nil)

	feedService := NewFeedService(ytM, ytdlM)
	f, _ := feedService.Create(context.Background(), "123")
	ti := time.Now()
	f.Content = podcast.New("title", "http://onet", "description", &ti, &ti)
	serialized, err := f.Serialize()

	assert.NoError(t, err, "serialize should not return error")
	assert.Contains(t, string(serialized), "<link>http://onet</link>")
}

func TestCreateShouldReturnFeedWithVideo1(t *testing.T) {
	ytM := new(YouTubeDependencyMock)
	ytM.On("ListEntry", context.Background(), "123").Return([]YouTubeFeedEntry{{
		ID:          "vi1",
		Title:       "vt1",
		Description: "vd1",
	}}, nil)
	ytM.On("GetChannelCache", context.Background(), "123").Return(YouTubeChannel{
		ChannelID:   "123",
		Country:     "pl",
		Description: "d1",
		Title:       "t1",
		URL:         "u1",
	}, nil)

	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetDetails", context.Background(), "vi1").Return(Details{}, nil)
	resp := httptest.NewRecorder()

	resp.Header().Set("Content-Type", "mp3")
	resp.Header().Set("Content-Length", "123")

	feedService := NewFeedService(ytM, ytdlM)
	f, _ := feedService.Create(context.Background(), "123")

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
	ytdlM.On("GetDetails", context.Background(), id).Return(Details{
		URL: *u,
	}, nil)

	item := FeedItem{ID: id}
	expectedItem := FeedItem{
		ID:         id,
		FileURL:    u,
		Duration:   0,
		FileLength: 0,
	}
	feedService := NewFeedService(YouTubeService{}, ytdlM)

	feedItem, err := feedService.Enrich(context.Background(), item)
	assert.NoError(t, err)
	assert.EqualValues(t, feedItem, expectedItem)
}

func TestItem_enrichReturnError(t *testing.T) {
	errorMessage := "error message"
	id := "123"
	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetDetails", context.Background(), id).Return(Details{}, errors.New(errorMessage))

	item := FeedItem{ID: id}
	feedService := NewFeedService(YouTubeService{}, ytdlM)

	_, err := feedService.Enrich(context.Background(), item)
	assert.Error(t, err, err)
}

func TestFeed_GetDetails(t *testing.T) {
	ytM := new(YouTubeDependencyMock)
	file := NewFileService()

	feedService := NewFeedService(ytM, file)

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
				ytM.On("GetChannelCache", context.Background(), "1").Return(YouTubeChannel{
					ChannelID:   "ch1",
					Country:     "pl1",
					Description: "d1",
					PublishedAt: time.Date(2000, 11, 01, 1, 1, 1, 1, time.UTC),
					Thumbnail:   YouTubeThumbnail{},
					Title:       "t1",
					URL:         "u1",
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
				ytM.On("GetChannelCache", context.Background(), "2").Return(YouTubeChannel{}, errors.New("can't get channel"))
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

			pod, err := feedService.GetFeedInformation(context.Background(), tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedInformation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if pod.Title != tt.expect.Title &&
				pod.Description != tt.expect.Description &&
				pod.Link != tt.expect.Link {
				t.Errorf("GetFeedInformation() get = %v, want = %v", f.Content, tt.expect)
			}
		})
	}
}

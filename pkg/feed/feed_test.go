package feed

import (
	"context"
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
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

func (m *YouTubeDependencyMock) ListEntry(s string) ([]app.YouTubeFeedEntry, error) {
	args := m.Called(s)
	return args.Get(0).([]app.YouTubeFeedEntry), args.Error(1)
}

func (m *YouTubeDependencyMock) GetChannelCache(s string) (app.YouTubeChannel, error) {
	args := m.Called(s)
	return args.Get(0).(app.YouTubeChannel), args.Error(1)
}

type YTDLDependencyMock struct {
	mock.Mock
}

func (m *YTDLDependencyMock) GetFileURL(info *ytdl.VideoInfo) (url.URL, error) {
	args := m.Called(info)
	return args.Get(0).(url.URL), args.Error(1)
}

func (m *YTDLDependencyMock) GetFileInformation(videoId string) (app.YTDLVideo, error) {
	args := m.Called(videoId)
	return args.Get(0).(app.YTDLVideo), args.Error(1)
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

func TestFeed_serialize(t *testing.T) {
	d := new(YouTubeDependencyMock)
	d.On("ListEntry", "123").Return([]app.YouTubeFeedEntry{}, nil)
	d.On("GetChannelCache", "123").Return(app.YouTubeChannel{}, nil)

	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileURL", "123").Return([]app.YouTubeFeedEntry{}, nil)
	ytdlM.On("GetFileInformation", "123").Return(app.YouTubeChannel{}, nil)

	f, _ := Create("123", Dependencies{YouTube: d, YTDL: ytdlM})
	ti := time.Now()
	f.Content = podcast.New("title", "http://onet", "description", &ti, &ti)
	serialized, err := f.Serialize()

	assert.NoError(t, err, "serialize should not return error")
	assert.Contains(t, string(serialized), "<link>http://onet</link>")
}

func TestCreateShouldReturnFeedWithVideo1(t *testing.T) {
	d := new(YouTubeDependencyMock)
	d.On("ListEntry", "123").Return([]app.YouTubeFeedEntry{{
		ID:          "vi1",
		Title:       "vt1",
		Description: "vd1",
	}}, nil)
	d.On("GetChannelCache", "123").Return(app.YouTubeChannel{
		ChannelId:   "123",
		Country:     "pl",
		Description: "d1",
		Title:       "t1",
		Url:         "u1",
	}, nil)

	vd := new(DependencyVideoMock)
	fileURLRaw := "http://onet.pl"
	fileURL, _ := url.Parse(fileURLRaw)
	info := app.YTDLVideo{
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

	f, _ := Create("123", Dependencies{
		YouTube: d,
		YTDL:    ytdlM,
	})

	assert.Equal(t, "123", f.ChannelID)
	assert.Len(t, f.Content.Items, 1)
	assert.Equal(t, "vt1", f.Items[0].Title)
	assert.Equal(t, "vd1", f.Items[0].Description)
	assert.True(t, d.AssertNumberOfCalls(t, "ListEntry", 1))
	assert.True(t, d.AssertNumberOfCalls(t, "GetChannelCache", 1))
	assert.Equal(t, "pl", f.Content.Language)
	assert.Equal(t, "d1", f.Content.Description)
	assert.Equal(t, "t1", f.Content.Title)
	assert.Equal(t, "u1", f.Content.Link)
}

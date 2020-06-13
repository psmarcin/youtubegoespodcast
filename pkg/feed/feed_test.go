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
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type DependencyMock struct {
	mock.Mock
}

func (m *DependencyMock) VideosList(s string) ([]youtube.Video, error) {
	args := m.Called(s)
	return args.Get(0).([]youtube.Video), args.Error(1)
}

func (m *DependencyMock) ChannelsGetFromCache(s string) (youtube.Channel, error) {
	args := m.Called(s)
	return args.Get(0).(youtube.Channel), args.Error(1)
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
	d := new(DependencyMock)
	d.On("VideosList", "123").Return([]youtube.Video{}, nil)
	d.On("ChannelsGetFromCache", "123").Return(youtube.Channel{}, nil)

	f, _ := Create("123", Dependencies{YouTube: d})
	ti := time.Now()
	f.Content = podcast.New("title", "http://onet", "description", &ti, &ti)
	serialized, err := f.Serialize()

	assert.NoError(t, err, "serialize should not return error")
	assert.Contains(t, string(serialized), "<link>http://onet</link>")
}

func TestCreateShouldReturnFeedWithVideo1(t *testing.T) {
	d := new(DependencyMock)
	d.On("VideosList", "123").Return([]youtube.Video{{
		ID:          "vi1",
		Title:       "vt1",
		Description: "vd1",
		Url:         "vu1",
	}}, nil)
	d.On("ChannelsGetFromCache", "123").Return(youtube.Channel{
		ChannelId:   "123",
		Country:     "pl",
		Description: "d1",
		Title:       "t1",
		Url:         "u1",
	}, nil)

	vd := new(DependencyVideoMock)
	fileURLRaw := "http://onet.pl"
	fileURL, _ := url.Parse(fileURLRaw)
	info := &ytdl.VideoInfo{
		ID:              "",
		Title:           "Title Video",
		Description:     "Description Video",
		DatePublished:   time.Time{},
		Formats:         nil,
		DASHManifestURL: "",
		HLSManifestURL:  "",
		Keywords:        nil,
		Uploader:        "",
		Song:            "",
		Artist:          "",
		Album:           "",
		Writers:         "",
		Duration:        0,
	}
	vd.On("GetFileUrl", info).Return(fileURL, nil)
	vd.On("Info", context.Background(), "vi1").Return(info, nil)
	resp := httptest.NewRecorder()
	resp.Header().Set("Content-Type", "mp3")
	resp.Header().Set("Content-Length", "123")
	vd.On("Details", fileURLRaw).Return(resp.Result(), nil)

	f, _ := Create("123", Dependencies{
		YouTube: d,
		Video:   video.Dependencies{Details: vd.Details, Info: vd.Info, GetFileUrl: vd.GetFileUrl},
	})

	assert.Equal(t, f.ChannelID, "123")
	assert.Len(t, f.Content.Items, 1)
	assert.Equal(t, f.Items[0].Details.ContentLength, int64(123))
	assert.Equal(t, f.Items[0].Details.ContentType, "mp3")
	assert.Equal(t, f.Items[0].Details.Title, "Title Video")
	assert.Equal(t, f.Items[0].Details.Description, "Description Video")
	assert.Equal(t, f.Items[0].Details.FileUrl.String(), "http://onet.pl")
	assert.True(t, d.AssertNumberOfCalls(t, "VideosList", 1))
	assert.True(t, d.AssertNumberOfCalls(t, "ChannelsGetFromCache", 1))
	assert.Equal(t, f.Content.Language, "pl")
	assert.Equal(t, f.Content.Description, "d1")
	assert.Equal(t, f.Content.Title, "t1")
	assert.Equal(t, f.Content.Link, "u1")
}

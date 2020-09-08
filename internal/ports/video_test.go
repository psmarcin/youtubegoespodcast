package ports

import (
	application "github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/rylio/ytdl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

type YTDLDependencyMock struct {
	mock.Mock
}

func (m *YTDLDependencyMock) GetFileURL(info *ytdl.VideoInfo) (url.URL, error) {
	args := m.Called(info)
	return args.Get(0).(url.URL), args.Error(1)
}

func (m *YTDLDependencyMock) GetFileInformation(videoId string) (application.YTDLVideo, error) {
	args := m.Called(videoId)
	return args.Get(0).(application.YTDLVideo), args.Error(1)
}

func TestHandler(t *testing.T) {
	type args struct {
		r *http.Request
	}

	req, _ := http.NewRequest(http.MethodGet, "/video/ulCdoCfw-bY/track.mp3", nil)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "videoID: ulCdoCfw-bY",
			args: args{
				r: req,
			},
		},
	}
	u1, _ := url.Parse("http://yt.com")
	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileInformation", "ulCdoCfw-bY").Return(application.YTDLVideo{
		ID:          "JZAunPKoHL0",
		URL:         u1,
		Description: "d2",
		Title:       "t1",
		FileUrl:     u1,
	}, nil)

	fiberServer := CreateHTTPServer()
	h := NewHttpServer(fiberServer, application.YouTubeService{}, ytdlM)
	app := h.Serve()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := app.Test(tt.args.r, -1)
			if err != nil {
				t.Errorf("should not throw error on app start")
			}
			assert.Equal(t, http.StatusFound, resp.StatusCode)
		})
	}
}

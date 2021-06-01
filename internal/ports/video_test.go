package ports

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	application "github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type YTDLDependencyMock struct {
	mock.Mock
}

func (m *YTDLDependencyMock) GetDetails(ctx context.Context, videoID string) (application.Details, error) {
	args := m.Called(mock.Anything, videoID)
	return args.Get(0).(application.Details), args.Error(1)
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
	u1, _ := url.Parse("http://youtube.com")
	ytdlM := new(YTDLDependencyMock)
	ytdlM.On("GetFileInformation", context.Background(), "ulCdoCfw-bY").Return(application.Details{
		URL: *u1,
	}, nil)

	fiberServer := CreateHTTPServer()
	h := NewHTTPServer(fiberServer, application.YouTubeService{}, application.NewFileService())
	app := h.Serve()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := app.Test(tt.args.r, -1)
			if err != nil {
				t.Errorf("should not throw error on app start")
			}
			defer resp.Body.Close()
			assert.Equal(t, http.StatusFound, resp.StatusCode)
		})
	}
}

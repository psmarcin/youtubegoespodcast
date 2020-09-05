package api

import (
	application "github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

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

	deps := Dependencies{
		YouTube: application.YouTubeService{},
	}
	app := Start(deps)

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

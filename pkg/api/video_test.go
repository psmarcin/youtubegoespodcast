package api

import (
	"github.com/psmarcin/youtubegoespodcast/pkg/cache"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
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

	c, _ := cache.Connect()

	deps := Dependencies{
		Cache:   c,
		YouTube: youtube.YT{},
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

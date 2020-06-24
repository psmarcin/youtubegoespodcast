package api

import (
	"github.com/psmarcin/youtubegoespodcast/pkg/cache"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_errorHandler_500(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, "/asdas", nil)

	c, _ := cache.Connect()

	deps := Dependencies{
		Cache:   c,
		YouTube: youtube.YT{},
	}
	app := Start(deps)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Errorf("should not throw error on app start")
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

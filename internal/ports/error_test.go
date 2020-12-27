package ports

import (
	application "github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_errorHandler_500(t *testing.T) {
	req, _ := http.NewRequest(http.MethodHead, "/asdas", nil)

	fiberServer := CreateHTTPServer()
	h := NewHttpServer(fiberServer, application.YouTubeService{}, application.NewFileService())

	app := h.Serve()

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Errorf("should not throw error on app start")
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

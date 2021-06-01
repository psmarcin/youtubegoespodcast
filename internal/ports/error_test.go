package ports

import (
	"context"
	"net/http"
	"testing"

	application "github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/stretchr/testify/assert"
)

func Test_errorHandler_500(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodHead, "/asdas", nil)
	fiberServer := CreateHTTPServer()
	h := NewHTTPServer(fiberServer, application.YouTubeService{}, application.NewFileService())

	app := h.Serve()

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Errorf("should not throw error on app start")
	}
	defer resp.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

package trending

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		args   args
		status int
	}{
		{
			name: "status code OK, empty content-type and not empty body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/target", nil),
			},
			status: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution
			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/videos").
				Reply(http.StatusOK).
				BodyString("{}")

			Handler(tt.args.w, tt.args.r)
			assert.Equal(t, tt.args.w.Code, tt.status)
			assert.NotEqual(t, tt.args.w.HeaderMap.Get("content-type"), "application/json")
			assert.NotEmpty(t, tt.args.w.Body)
		})
	}
}

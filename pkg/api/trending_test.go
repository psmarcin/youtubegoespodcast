package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestTrendingHandler(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name   string
		args   args
		status string
	}{
		{
			name: "status code OK, empty content-type and not empty body",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/target", nil),
			},
			status: "404 Not Found",
		},
	}

	app := Start()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution
			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/videos").
				Reply(http.StatusOK).
				BodyString("{}")

			resp, err := app.Test(tt.args.r)
			if err != nil {
				t.Errorf("should not throw error on app start")
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("should not throw error on body serialization")
			}

			assert.Equal(t, tt.status, resp.Status)
			assert.NotEqual(t, resp.Header.Get("content-type"), "application/json")
			assert.NotEmpty(t, body)
		})
	}
}

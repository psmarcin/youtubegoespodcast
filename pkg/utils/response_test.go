package utils

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONResponse(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should set content-type to application/json",
			args: args{
				w: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			JSONResponse(tt.args.w)
			w := tt.args.w.Result()
			assert.Equal(w.Header.Get("content-type"), "application/json")
		})
	}
}

func TestOkResponse(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		w *httptest.ResponseRecorder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should set status code",
			args: args{
				w: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OkResponse(tt.args.w)
			w := tt.args.w.Result()
			assert.Equal(w.StatusCode, http.StatusOK)
		})
	}
}

func TestWriteBodyResponse(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		w    *httptest.ResponseRecorder
		body string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should write body",
			args: args{
				w:    httptest.NewRecorder(),
				body: "testBody",
			},
		},
		{
			name: "Should write body",
			args: args{
				w:    httptest.NewRecorder(),
				body: "askdk asldlk asjdlj123l'4jhi 1;lke🤙",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteBodyResponse(tt.args.w, tt.args.body)
			w := tt.args.w.Result()
			body, _ := ioutil.ReadAll(w.Body)
			assert.Equal(string(body), tt.args.body)
		})
	}
}

func TestAllowCorsResponse(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	type args struct {
		w      *httptest.ResponseRecorder
		r      *http.Request
		origin string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Set proper CORS headers",
			args: args{
				w:      httptest.NewRecorder(),
				r:      req,
				origin: "http://localhost:8080",
			},
			want: "http://localhost:8080",
		},
		{
			name: "Don't set CORS headers for not allowe domain",
			args: args{
				w:      httptest.NewRecorder(),
				r:      req,
				origin: "http://fake.domain",
			},
			want: "",
		},
		{
			name: "Set CORS headers for allowed domain https://yt.psmarcin.dev",
			args: args{
				w:      httptest.NewRecorder(),
				r:      req,
				origin: "https://yt.psmarcin.dev",
			},
			want: "https://yt.psmarcin.dev",
		},
		{
			name: "Set CORS headers for allowed domain https://yt.psmarcin.me",
			args: args{
				w:      httptest.NewRecorder(),
				r:      req,
				origin: "https://yt.psmarcin.me",
			},
			want: "https://yt.psmarcin.me",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Header.Set("origin", tt.args.origin)
			AllowCorsResponse(tt.args.w, tt.args.r)
			response := tt.args.w.Result()
			assert.Equal(t, response.Header.Get("Access-Control-Allow-Origin"), tt.want)
		})
	}
}

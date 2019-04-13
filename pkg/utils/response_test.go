package utils

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

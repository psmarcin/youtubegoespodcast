package video

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	req := httptest.NewRequest(http.MethodGet, "/video/videoId.mp3", nil)
	req = mux.SetURLVars(req, map[string]string{"videoId": "NtQkz0aRDe8"})

	tests := []struct {
		name string
		args args
	}{
		{
			name: "videoID: uwfdFCP3KYM",
			args: args{
				w: httptest.NewRecorder(),
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Handler(tt.args.w, tt.args.r)
			if tt.args.w.Code != http.StatusFound {
				t.Errorf("Handler(): %+v want %+v ", tt.args.w.Code, http.StatusFound)
			}
		})
	}
}

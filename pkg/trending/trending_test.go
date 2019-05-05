package trending

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "status code OK, empty content-type and not empty body",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/target", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Handler(tt.args.w, tt.args.r)
			assert.Equal(t, tt.args.w.Code, http.StatusOK)
			assert.NotEqual(t, tt.args.w.HeaderMap.Get("content-type"), "application/json")
			assert.NotEmpty(t, tt.args.w.Body)
		})
	}
}

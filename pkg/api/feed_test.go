package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handleError(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want struct {
			code int
		}
	}{
		{
			name: "Write error",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
		},
	}
	app := Start()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := app.Test(tt.args.r)
			if err != nil {
				t.Errorf("should not throw error on app start")
			}

			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

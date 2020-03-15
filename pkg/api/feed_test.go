package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handleError(t *testing.T) {
	type args struct {
		w   *httptest.ResponseRecorder
		err error
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
				w:   httptest.NewRecorder(),
				err: errors.New("Bad request"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleError(tt.args.w, tt.args.err)
			assert.Equal(t, tt.args.w.Code, http.StatusBadRequest)
			assert.Equal(t, tt.args.w.Body.String(), tt.args.err.Error())
		})
	}
}

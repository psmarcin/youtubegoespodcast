package youtube

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRequest(t *testing.T) {
	type json struct {
		Test string `json:"argument"`
	}
	type args struct {
		req *http.Request
		y   interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "request to trending",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/test", nil),
				y:   json{},
			},
		},
	}
	for _, tt := range tests {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("https://www.googleapis.com").
			Get("/test").
			Reply(200).
			BodyString(`{
					"argument": "123"
					}
				`)
		t.Run(tt.name, func(t *testing.T) {
			Request(tt.args.req, tt.args.y)
			assert.Equal(t, tt.args.y, "123")
		})
	}
}

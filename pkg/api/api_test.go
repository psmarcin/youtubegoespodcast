package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_rootHandler(t *testing.T) {
	t.Run("respose ok", func(t *testing.T) {
		assert.HTTPSuccess(t, rootHandler, http.MethodGet, "/", nil, nil)
		assert.HTTPBodyContains(t, rootHandler, http.MethodGet, "/", nil, "{\"status\": \"OK\"}")
	})
}

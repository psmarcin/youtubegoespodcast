package api

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func Test_rootHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	app := Start()
	resp, err := app.Test(req)
	assert.Empty(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Empty(t, err)

	assert.True(t, strings.Contains(string(body), "{\"status\":true"))
}

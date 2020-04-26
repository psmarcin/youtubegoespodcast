package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/h2non/gock.v1"
)

func TestChannelsHandler(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type mock struct {
		statusCode int
		body       string
	}
	type expect struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		args   args
		mock   mock
		expect expect
	}{
		{
			name: "should return collection of channels",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/channels", nil),
			},
			mock: mock{
				statusCode: http.StatusOK,
				body:       successResponseIzak,
			},
			expect: expect{
				statusCode: http.StatusOK,
				body:       `[{"channelId":"UCSeqSPSJDl0EXezzjJrwQNg","thumbnail":"https://yt3.ggpht.com/-oRBv14_mRi0/AAAAAAAAAAI/AAAAAAAAAAA/iAsVn9M_pzY/s800-c-k-no-mo-rj-c0xffffff/photo.jpg","title":"izak LIVE"},{"channelId":"UCKIJxLhLvwlna82NctN6Zow","thumbnail":"https://yt3.ggpht.com/-KwxwQ9zLUvM/AAAAAAAAAAI/AAAAAAAAAAA/uKA1BY19Cik/s800-c-k-no-mo-rj-c0xffffff/photo.jpg","title":"Izak"},{"channelId":"UCNEjf2pgf4-9wlp1E45-4Yw","thumbnail":"https://yt3.ggpht.com/-9C0hkkXpnmA/AAAAAAAAAAI/AAAAAAAAAAA/y9ZMgBpkXSI/s800-c-k-no-mo-rj-c0xffffff/photo.jpg","title":"izak LIVE"},{"channelId":"UCtgwr2O1Ti8LzfDr3rBPONw","thumbnail":"https://yt3.ggpht.com/-gw-dgVl5yE4/AAAAAAAAAAI/AAAAAAAAAAA/bw2Dhbm5hvc/s800-c-k-no-mo-rj-c0xffffff/photo.jpg","title":"Izak"},{"channelId":"UCemUDaS2uHy5jxyoo_AopmA","thumbnail":"https://yt3.ggpht.com/-OirIp9BLN7g/AAAAAAAAAAI/AAAAAAAAAAA/yqgL4OrOXpY/s800-c-k-no-mo-rj-c0xffffff/photo.jpg","title":"IZAK 10"}]`,
			},
		},
		{
			name: "should return 400 error on youtube 400",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/channels", nil),
			},
			mock: mock{
				statusCode: http.StatusBadRequest,
				body:       "",
			},
			expect: expect{
				statusCode: http.StatusNotFound,
				body:       "Not Found",
			},
		},
	}

	app := Start()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution
			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/search").
				Reply(tt.mock.statusCode).
				BodyString(tt.mock.body)

			resp, _ := app.Test(tt.args.r)
			body, _ := ioutil.ReadAll(resp.Body)

			assert.Equal(t, tt.expect.statusCode, resp.StatusCode)
			assert.Equal(t, tt.expect.body, string(body))
		})
	}
}

const (
	successResponseIzak = `{
		"kind": "youtube#searchListResponse",
		"etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/fLpK0b3wYvlXpqK68_8oJdeH6bM\"",
		"nextPageToken": "CAUQAA",
		"regionCode": "PL",
		"pageInfo": {
			"totalResults": 12762,
			"resultsPerPage": 5
		},
		"items": [
			{
			"kind": "youtube#searchResult",
			"etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/1O1MtuvMO2UyPSBhsSTY_HZemK8\"",
			"id": {
				"kind": "youtube#channel",
				"channelId": "UCSeqSPSJDl0EXezzjJrwQNg"
			},
			"snippet": {
				"publishedAt": "2006-12-01T09:59:48.000Z",
				"channelId": "UCSeqSPSJDl0EXezzjJrwQNg",
				"title": "izak LIVE",
				"description": "Oficjalny kanał Piotra 'izaka' Skowyrskiego.",
				"thumbnails": {
				"default": {
					"url": "https://yt3.ggpht.com/-oRBv14_mRi0/AAAAAAAAAAI/AAAAAAAAAAA/iAsVn9M_pzY/s88-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"medium": {
					"url": "https://yt3.ggpht.com/-oRBv14_mRi0/AAAAAAAAAAI/AAAAAAAAAAA/iAsVn9M_pzY/s240-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"high": {
					"url": "https://yt3.ggpht.com/-oRBv14_mRi0/AAAAAAAAAAI/AAAAAAAAAAA/iAsVn9M_pzY/s800-c-k-no-mo-rj-c0xffffff/photo.jpg"
				}
				},
				"channelTitle": "izak LIVE",
				"liveBroadcastContent": "upcoming"
			}
			},
			{
			"kind": "youtube#searchResult",
			"etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/D8vFHsAVO8qB-KfgDCXVIdsPyqU\"",
			"id": {
				"kind": "youtube#channel",
				"channelId": "UCKIJxLhLvwlna82NctN6Zow"
			},
			"snippet": {
				"publishedAt": "2016-01-30T09:29:32.000Z",
				"channelId": "UCKIJxLhLvwlna82NctN6Zow",
				"title": "Izak",
				"description": "Fuck à hoe.",
				"thumbnails": {
				"default": {
					"url": "https://yt3.ggpht.com/-KwxwQ9zLUvM/AAAAAAAAAAI/AAAAAAAAAAA/uKA1BY19Cik/s88-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"medium": {
					"url": "https://yt3.ggpht.com/-KwxwQ9zLUvM/AAAAAAAAAAI/AAAAAAAAAAA/uKA1BY19Cik/s240-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"high": {
					"url": "https://yt3.ggpht.com/-KwxwQ9zLUvM/AAAAAAAAAAI/AAAAAAAAAAA/uKA1BY19Cik/s800-c-k-no-mo-rj-c0xffffff/photo.jpg"
				}
				},
				"channelTitle": "Izak",
				"liveBroadcastContent": "none"
			}
			},
			{
			"kind": "youtube#searchResult",
			"etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/jGShR8pvvEsyo1ciKF_6_GdiYts\"",
			"id": {
				"kind": "youtube#channel",
				"channelId": "UCNEjf2pgf4-9wlp1E45-4Yw"
			},
			"snippet": {
				"publishedAt": "2014-09-22T23:33:58.000Z",
				"channelId": "UCNEjf2pgf4-9wlp1E45-4Yw",
				"title": "izak LIVE",
				"description": "No to jak gramy w WOTA Zapraszam was do tej gry jest ona bardzo fajna :) A oto link do niej:http://worldoftanks.eu/ Loguj się i graj razem z ProB'em :) Wbij tak, ...",
				"thumbnails": {
				"default": {
					"url": "https://yt3.ggpht.com/-9C0hkkXpnmA/AAAAAAAAAAI/AAAAAAAAAAA/y9ZMgBpkXSI/s88-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"medium": {
					"url": "https://yt3.ggpht.com/-9C0hkkXpnmA/AAAAAAAAAAI/AAAAAAAAAAA/y9ZMgBpkXSI/s240-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"high": {
					"url": "https://yt3.ggpht.com/-9C0hkkXpnmA/AAAAAAAAAAI/AAAAAAAAAAA/y9ZMgBpkXSI/s800-c-k-no-mo-rj-c0xffffff/photo.jpg"
				}
				},
				"channelTitle": "izak LIVE",
				"liveBroadcastContent": "none"
			}
			},
			{
			"kind": "youtube#searchResult",
			"etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/NlqUAj4IPYmWq864S6qcgqTqNIg\"",
			"id": {
				"kind": "youtube#channel",
				"channelId": "UCtgwr2O1Ti8LzfDr3rBPONw"
			},
			"snippet": {
				"publishedAt": "2015-05-01T10:04:44.000Z",
				"channelId": "UCtgwr2O1Ti8LzfDr3rBPONw",
				"title": "Izak",
				"description": "I hope to make this account a place where I can post guitar covers, tutorials, and originals. It would be a lot of fun to learn from others and maybe eventually ...",
				"thumbnails": {
				"default": {
					"url": "https://yt3.ggpht.com/-gw-dgVl5yE4/AAAAAAAAAAI/AAAAAAAAAAA/bw2Dhbm5hvc/s88-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"medium": {
					"url": "https://yt3.ggpht.com/-gw-dgVl5yE4/AAAAAAAAAAI/AAAAAAAAAAA/bw2Dhbm5hvc/s240-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"high": {
					"url": "https://yt3.ggpht.com/-gw-dgVl5yE4/AAAAAAAAAAI/AAAAAAAAAAA/bw2Dhbm5hvc/s800-c-k-no-mo-rj-c0xffffff/photo.jpg"
				}
				},
				"channelTitle": "Izak",
				"liveBroadcastContent": "none"
			}
			},
			{
			"kind": "youtube#searchResult",
			"etag": "\"Bdx4f4ps3xCOOo1WZ91nTLkRZ_c/DStm85L9Ik4UeJH2D0FUV_N01Yg\"",
			"id": {
				"kind": "youtube#channel",
				"channelId": "UCemUDaS2uHy5jxyoo_AopmA"
			},
			"snippet": {
				"publishedAt": "2017-02-18T06:54:05.000Z",
				"channelId": "UCemUDaS2uHy5jxyoo_AopmA",
				"title": "IZAK 10",
				"description": "Carevi dobrodosli na moj kanal zovem se izak i udite u moj klan na brawl starsu zove se Zaki armija.",
				"thumbnails": {
				"default": {
					"url": "https://yt3.ggpht.com/-OirIp9BLN7g/AAAAAAAAAAI/AAAAAAAAAAA/yqgL4OrOXpY/s88-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"medium": {
					"url": "https://yt3.ggpht.com/-OirIp9BLN7g/AAAAAAAAAAI/AAAAAAAAAAA/yqgL4OrOXpY/s240-c-k-no-mo-rj-c0xffffff/photo.jpg"
				},
				"high": {
					"url": "https://yt3.ggpht.com/-OirIp9BLN7g/AAAAAAAAAAI/AAAAAAAAAAA/yqgL4OrOXpY/s800-c-k-no-mo-rj-c0xffffff/photo.jpg"
				}
				},
				"channelTitle": "IZAK 10",
				"liveBroadcastContent": "upcoming"
			}
			}
		]
		}`
)

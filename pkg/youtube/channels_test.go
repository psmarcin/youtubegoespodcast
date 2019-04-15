package youtube

import (
	"gopkg.in/h2non/gock.v1"
	"reflect"
	"testing"
)

func TestGetChannels(t *testing.T) {
	type args struct {
		q string
	}
	tests := []struct {
		name string
		args args
		want YoutubeResponse
	}{
		{
			name: "no q query param",
			args: args{
				q: "",
			},
			want: YoutubeResponse{
				Items: []Items{Items{
					Snippet: Snippet{
						ChannelID:    "UCqrBKQHQVEvD2Q9FP0DCP2g",
						ChannelTitle: "benny blanco",
						Title:        "benny blanco, Halsey &amp; Khalid – Eastside (official video)",
						Thumbnails: Thumbnails{
							High: High{
								URL:    "https://i.ytimg.com/vi/56WBK4ZK_cw/hqdefault.jpg",
								Width:  480,
								Height: 360,
							},
						},
					},
				},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off() // Flush pending mocks after test execution

			gock.New("https://www.googleapis.com").
				Get("/youtube/v3/search").
				Reply(200).
				BodyString(`{
					"items": [
						{
						"snippet": {
							"channelId": "UCqrBKQHQVEvD2Q9FP0DCP2g",
							"title": "benny blanco, Halsey &amp; Khalid – Eastside (official video)",
							"description": "\"Eastside\" out now: http://smarturl.it/EastsideBB Directed by Jake Schreier Executive Producers: Jackie Kelman Bisbee, Alex Fisch Line Producer: Molly Gale ...",
							"thumbnails": {
							"high": {
								"url": "https://i.ytimg.com/vi/56WBK4ZK_cw/hqdefault.jpg",
								"width": 480,
								"height": 360
							}
							},
							"channelTitle": "benny blanco"
						}
						}
					]
					}
				`)
			if got := GetChannels(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChannels() = %v, want %v", got, tt.want)
			}
		})
	}
}

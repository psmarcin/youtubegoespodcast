package youtube

import "testing"

func TestSerialize(t *testing.T) {
	type args struct {
		channelResponse YoutubeResponse
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "serialize empty",
			args: args{
				channelResponse: YoutubeResponse{},
			},
			want: "[]",
		},
		{
			name: "should serialize youtube response",
			args: args{
				channelResponse: YoutubeResponse{
					Items: []Items{
						{
							Snippet: Snippet{
								ChannelID:    "ChannelID",
								ChannelTitle: "ChannelTitle",
								Thumbnails: Thumbnails{
									High: High{
										Height: 100,
										URL:    "http://google.go",
										Width:  200,
									},
								},
								Title: "Title",
							},
						},
					},
				},
			},
			want: "[{\"channelId\":\"ChannelID\",\"thumbnail\":\"http://google.go\",\"title\":\"ChannelTitle\"}]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Serialize(tt.args.channelResponse); got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

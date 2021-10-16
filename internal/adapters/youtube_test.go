package adapters

import (
	"net/url"
	"testing"
	"time"

	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/youtube/v3"
)

func Test_mapChannelItemToYouTubeChannels(t *testing.T) {
	pubString := "2006-01-02T15:04:05Z"
	pubTime, _ := time.Parse(time.RFC3339, pubString)
	urlString := "https://yt.google.com"
	urlParsed, _ := url.Parse(urlString)

	type args struct {
		item *youtube.Channel
	}
	tests := []struct {
		name    string
		args    args
		want    app.YouTubeChannel
		wantErr bool
	}{
		{
			name: "should map thumbnail",
			args: args{
				item: &youtube.Channel{
					Snippet: &youtube.ChannelSnippet{
						Thumbnails: &youtube.ThumbnailDetails{
							High: &youtube.Thumbnail{
								Height: 500,
								Url:    urlString,
								Width:  300,
							},
						},
						PublishedAt: "2006-01-02T15:04:05Z",
					},
					TopicDetails: &youtube.ChannelTopicDetails{
						TopicCategories: []string{},
					},
				},
			},
			want: app.YouTubeChannel{
				AuthorEmail: "email@example.com",
				PublishedAt: pubTime,
				Thumbnail: app.YouTubeThumbnail{
					Height: 500,
					Width:  300,
					URL:    *urlParsed,
				},
				URL: YouTubeChannelBaseURL,
			},
			wantErr: false,
		},
		{
			name: "should map thumbnail even without valid url",
			args: args{
				item: &youtube.Channel{
					Snippet: &youtube.ChannelSnippet{
						Thumbnails: &youtube.ThumbnailDetails{
							High: &youtube.Thumbnail{
								Height: 500,
								Url:    "x",
								Width:  300,
							},
						},
						PublishedAt: "2006-01-02T15:04:05Z",
					},
					TopicDetails: &youtube.ChannelTopicDetails{
						TopicCategories: []string{},
					},
				},
			},
			want: app.YouTubeChannel{
				AuthorEmail: "email@example.com",
				PublishedAt: pubTime,
				Thumbnail: app.YouTubeThumbnail{
					Height: 500,
					Width:  300,
					URL: url.URL{
						Path: "x",
					},
				},
				URL: YouTubeChannelBaseURL,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapChannelToYouTubeChannel(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapChannelToYouTubeChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func Test_mapSearchItemToYouTubeChannel(t *testing.T) {
	pubString := "2006-01-02T15:04:05Z"
	pubTime, _ := time.Parse(time.RFC3339, pubString)
	urlString := "https://yt.google.com"
	urlParsed, _ := url.Parse(urlString)

	type args struct {
		item *youtube.SearchResult
	}
	tests := []struct {
		name    string
		args    args
		want    app.YouTubeChannel
		wantErr bool
	}{
		{
			name: "should map thumbnail",
			args: args{
				item: &youtube.SearchResult{
					Id: &youtube.ResourceId{
						ChannelId: "1",
					},
					Snippet: &youtube.SearchResultSnippet{
						PublishedAt: pubString,
						Thumbnails: &youtube.ThumbnailDetails{
							Default: &youtube.Thumbnail{
								Height: 500,
								Url:    urlString,
								Width:  1,
							},
						},
					},
				},
			},
			want: app.YouTubeChannel{
				ChannelID:   "1",
				PublishedAt: pubTime,
				Thumbnail: app.YouTubeThumbnail{
					Height: 500,
					Width:  1,
					URL:    *urlParsed,
				},
			},
			wantErr: false,
		},
		{
			name: "should map thumbnail url without valid url",
			args: args{
				item: &youtube.SearchResult{
					Id: &youtube.ResourceId{
						ChannelId: "1",
					},
					Snippet: &youtube.SearchResultSnippet{
						PublishedAt: pubString,
						Thumbnails: &youtube.ThumbnailDetails{
							Default: &youtube.Thumbnail{
								Height: 500,
								Url:    "x",
								Width:  1,
							},
						},
					},
				},
			},
			want: app.YouTubeChannel{
				ChannelID:   "1",
				PublishedAt: pubTime,
				Thumbnail: app.YouTubeThumbnail{
					Height: 500,
					Width:  1,
					URL: url.URL{
						Path: "x",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapSearchItemToYouTubeChannel(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapSearchItemToYouTubeChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValues(t, tt.want, got)
		})
	}
}

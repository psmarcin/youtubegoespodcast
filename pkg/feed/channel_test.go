package feed

import (
	"errors"
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"testing"
	"time"
)

type YT struct {
	returnChannel youtube.Channel
	returnError   error
}

func (yt YT) ChannelsGetFromCache(_ string) (youtube.Channel, error) {
	return yt.returnChannel, yt.returnError
}

func TestFeed_AddItem(t *testing.T) {
	type fields struct {
		ChannelID string
		Content   podcast.Podcast
	}
	type args struct {
		item podcast.Item
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should add new item",
			fields: fields{
				ChannelID: "123",
				Content:   podcast.Podcast{},
			},
			args: args{
				item: podcast.Item{
					Title:       "title123",
					Description: "descrioption123",
					GUID:        "123",
					Enclosure: &podcast.Enclosure{
						URL: "url",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{
				ChannelID: tt.fields.ChannelID,
				Content:   tt.fields.Content,
			}
			if err := f.AddItem(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeed_AddItem_CountAddedItems(t *testing.T) {
	f := &Feed{}
	// should add new item
	err := f.AddItem(podcast.Item{
		Title:       "t1",
		Description: "d1",
		Enclosure: &podcast.Enclosure{
			URL: "u1",
		},
	})

	if err != nil {
		t.Errorf("AddItem() error = %v", err)
	}

	if len(f.Content.Items) != 1 {
		t.Errorf("AddItem() f.Item length = %v, want %v", len(f.Content.Items), 1)
	}

	// should not add new item (without title)
	err = f.AddItem(podcast.Item{
		Title: "",
		Enclosure: &podcast.Enclosure{
			URL: "u1",
		},
	})

	if len(f.Content.Items) != 1 {
		t.Errorf("AddItem() f.Item length = %v, want %v", len(f.Content.Items), 1)
	}

	// should not add new item (without url)
	err = f.AddItem(podcast.Item{
		Title: "t1",
		Enclosure: &podcast.Enclosure{
			URL: "",
		},
	})

	if len(f.Content.Items) != 1 {
		t.Errorf("AddItem() f.Item length = %v, want %v", len(f.Content.Items), 1)
	}
}

func TestFeed_GetDetails(t *testing.T) {
	type fields struct {
		ChannelID string
		Content   podcast.Podcast
	}
	type args struct {
		channelID string
		ch        YT
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		expect  podcast.Podcast
		wantErr bool
	}{
		{
			name: "should get channel details",
			fields: fields{
				ChannelID: "",
				Content:   podcast.Podcast{},
			}, args: args{
			ch: YT{
				returnChannel: youtube.Channel{
					ChannelId:   "ch1",
					Country:     "pl1",
					Description: "d1",
					PublishedAt: time.Date(2000, 11, 01, 1, 1, 1, 1, time.UTC),
					Thumbnail:   "th1",
					Title:       "t1",
					Url:         "u1",
				},
				returnError: nil,
			}},
			expect: podcast.Podcast{
				Title:       "t1",
				Link:        "u1",
				Description: "d1",
			},
			wantErr: false,
		},
		{
			name: "should throw error on youtube service error",
			fields: fields{
				ChannelID: "",
				Content:   podcast.Podcast{},
			}, args: args{
			ch: YT{
				returnChannel: youtube.Channel{},
				returnError:   errors.New("can't get channel"),
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Feed{
				ChannelID: tt.fields.ChannelID,
				Content:   tt.fields.Content,
			}

			err := f.GetDetails(tt.args.channelID, tt.args.ch)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetails() error = %v, wantErr %v", err, tt.wantErr)
			}

			if f.Content.Title != tt.expect.Title &&
				f.Content.Description != tt.expect.Description &&
				f.Content.Link != tt.expect.Link {
				t.Errorf("GetDetails() get = %v, want = %v", f.Content, tt.expect)
			}
		})
	}
}

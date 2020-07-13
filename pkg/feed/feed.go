package feed

import (
	"fmt"
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"time"
)

const (
	Generator = "YouTubeGoesPodcast/v2"
)

var l = logrus.WithField("source", "feed")

type Feed struct {
	ChannelID string
	Content   podcast.Podcast
	Items     []Item
}

type Dependencies struct {
	YouTube YouTubeDependency
}

type YouTubeDependency interface {
	VideosList(string) ([]youtube.Video, error)
	ChannelsGetFromCache(string) (youtube.Channel, error)
}

func Create(channelID string, dependencies Dependencies) (Feed, error) {
	f := Feed{
		ChannelID: channelID,
	}
	err := f.GetDetails(channelID, dependencies.YouTube)
	if err != nil {
		l.WithError(err).Errorf("can't get feed details for %s", channelID)
		return f, err
	}

	videos, err := dependencies.YouTube.VideosList(f.ChannelID)
	if err != nil {
		l.WithError(err).Errorf("can't get video list for %s", channelID)
		return f, err
	}

	err = f.SetItems(videos)
	if err != nil {
		l.WithError(err).Errorf("can't enrich videos")
		return f, err
	}

	err = f.SetVideos()
	if err != nil {
		l.WithError(err).Errorf("can't set videos for %s", channelID)
		return f, err
	}

	f.sortVideos()
	f.Content.Generator = Generator
	return f, nil
}

func (f *Feed) Serialize() ([]byte, error) {
	f.Content.Items = sortByOrder(f.Content.Items)
	return f.Content.Bytes(), nil
}

func (f *Feed) sortVideos() {
	f.Content.Items = sortByOrder(f.Content.Items)
}

func (f *Feed) SetVideos() error {
	// set channel last updated at field as latest item publishing date
	if len(f.Items) != 0 {
		lastItem := f.Items[0]
		f.Content.LastBuildDate = lastItem.PubDate.Format(time.RFC1123Z)
		f.Content.PubDate = f.Content.LastBuildDate
	}

	for _, item := range f.Items {
		videoURL := os.Getenv("API_URL") + "video/" + item.ID + "/track.mp3"

		err := f.AddItem(podcast.Item{
			GUID:        item.GUID,
			Title:       item.Title,
			Link:        item.Link.String(),
			Description: item.Description,
			PubDate:     &item.PubDate,
			Enclosure: &podcast.Enclosure{
				URL:    videoURL,
				Length: item.FileLength,
				Type:   podcast.MP3,
			},
			IDuration: fmt.Sprintf("%f", item.Duration.Seconds()),
			IExplicit: "no",
		})

		if err != nil {
			l.WithError(err).Errorf("can't add new video to feed")
			return err
		}
	}

	return nil
}

func (f *Feed) SetItems(videos []youtube.Video) error {
	items, err := NewMap(videos)
	if err != nil {
		return err
	}
	f.Items = items

	return nil
}

func sortByOrder(items []*podcast.Item) []*podcast.Item {
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(*items[j].PubDate)
	})
	return items
}

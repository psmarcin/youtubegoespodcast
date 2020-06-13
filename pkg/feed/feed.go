package feed

import (
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/video"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/sirupsen/logrus"
	"sort"
)

var l = logrus.WithField("source", "feed")

type Feed struct {
	ChannelID string
	Content   podcast.Podcast
	Items     []Item
}

type Dependencies struct {
	YouTube YouTubeDependency
	Video   video.Dependencies
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

	f.SetItems(videos)
	err = f.EnrichItems(dependencies.Video)
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
	f.Content.Generator = "xxx2"
	return f, nil
}

func (f *Feed) Serialize() ([]byte, error) {
	f.Content.Items = sortByOrder(f.Content.Items)
	return f.Content.Bytes(), nil
}

func (f *Feed) sortVideos() {
	f.Content.Items = sortByOrder(f.Content.Items)
}

func sortByOrder(items []*podcast.Item) []*podcast.Item {
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(*items[j].PubDate)
	})
	return items
}

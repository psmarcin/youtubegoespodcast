package feed

import (
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/sirupsen/logrus"
	"sort"
)

var l = logrus.WithField("source", "feed")

type Feed struct {
	ChannelID string
	Content   podcast.Podcast
}

func Create(channelID string) (Feed, error) {
	f := Feed{
		ChannelID: channelID,
	}
	err := f.GetDetails(channelID, youtube.Yt)
	if err != nil {
		l.WithError(err).Errorf("can't get feed details for %s", channelID)
		return f, err
	}

	videos, err := youtube.Yt.VideosList(f.ChannelID)
	if err != nil {
		l.WithError(err).Errorf("can't get video list for %s", channelID)
		return f, err
	}

	err = f.SetVideos(videos)
	if err != nil {
		l.WithError(err).Errorf("can't set videos for %s", channelID)
		return f, err
	}

	f.sortVideos()

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

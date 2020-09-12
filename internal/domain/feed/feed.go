package feed

import (
	"github.com/eduncan911/podcast"
	"github.com/sirupsen/logrus"
	"sort"
)

var l = logrus.WithField("source", "feed-domain")

func SortByPubDate(items []*podcast.Item) []*podcast.Item {
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(*items[j].PubDate)
	})
	return items
}

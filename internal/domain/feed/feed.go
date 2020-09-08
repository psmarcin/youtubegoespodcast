package feed

import (
	"github.com/eduncan911/podcast"
	"sort"
)

func SortByPubDate(items []*podcast.Item) []*podcast.Item {
	sort.Slice(items, func(i, j int) bool {
		return items[i].PubDate.After(*items[j].PubDate)
	})
	return items
}

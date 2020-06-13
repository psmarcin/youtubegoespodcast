package feed

import (
	"fmt"
	"github.com/eduncan911/podcast"
	"os"
	"time"

	vid "github.com/psmarcin/youtubegoespodcast/pkg/video"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
)

type Item struct {
	Video   youtube.Video
	Details vid.Details
}

func (f *Feed) EnrichItems(deps vid.Dependencies) error {
	stream := make(chan Item, len(f.Items))
	for i, video := range f.Items {
		go func(item Item, i int) error {
			vd, err := vid.GetDetails(item.Video.ID, deps)
			if err != nil {
				l.WithError(err).Error("item")
				stream <- Item{}
				return err
			}

			item.Details = vd

			stream <- item
			return nil
		}(video, i)
	}

	counter := 0
	var enriched []Item
	for {
		if counter >= len(f.Items) {
			break
		}
		enriched = append(enriched, <-stream)
		counter++
	}

	f.Items = enriched
	return nil
}

func (f *Feed) RemoveEmptyItems() {
	f.Items = filterEmptyVideoOut(f.Items)
}

func (f *Feed) SetVideos() error {
	// set channel last updated at field as latest item publishing date
	if len(f.Items) != 0 {
		lastItem := f.Items[0]
		f.Content.LastBuildDate = lastItem.Details.DatePublished.Format(time.RFC1123Z)
		f.Content.PubDate = f.Content.LastBuildDate
	}

	for _, item := range f.Items {
		videoURL := os.Getenv("API_URL") + "video/" + item.Video.ID + "/track.mp3"

		err := f.AddItem(podcast.Item{
			GUID:        item.Video.ID,
			Title:       item.Details.Title,
			Link:        item.Video.Url,
			Description: item.Details.Description,
			PubDate:     &item.Details.DatePublished,
			Enclosure: &podcast.Enclosure{
				URL:    videoURL,
				Length: item.Details.ContentLength,
				Type:   podcast.MP3,
			},
			IDuration: fmt.Sprintf("%f", item.Details.Duration.Seconds()),
			IExplicit: "no",
		})

		if err != nil {
			l.WithError(err).Errorf("can't add new video to feed")
			return err
		}
	}

	return nil
}

func (f *Feed) SetItems(videos []youtube.Video) {
	for _, video := range videos {
		f.Items = append(f.Items, Item{
			Video:   video,
			Details: vid.Details{},
		})
	}
}

func filterEmptyVideoOut(items []Item) []Item {
	var filtered []Item
	for _, item := range items {
		if item.Video.ID == "" {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered
}

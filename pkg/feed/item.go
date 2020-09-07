package feed

import (
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"net/url"
	"time"
)

type Item struct {
	ID          string
	GUID        string
	Title       string
	Link        *url.URL
	Description string
	PubDate     time.Time
	FileURL     *url.URL
	FileLength  int64
	FileType    string
	Duration    time.Duration
	IsExplicit  bool
}

type fileInformationGetter func(string) (app.YTDLVideo, error)

func NewMap(videos []app.YouTubeFeedEntry, dependencies YTDLDependencies) ([]Item, error) {
	stream := make(chan Item, len(videos))
	for _, v := range videos {
		go func(video app.YouTubeFeedEntry) {
			item, err := New(video, dependencies)
			if err != nil {
				l.WithError(err).WithField("video", video).Infof("can't create new item from video")
			}
			stream <- item
		}(v)
	}

	var items []Item
	counter := 0
	for {
		if counter >= len(videos) {
			break
		}

		item := <-stream
		if item.IsValid() {
			items = append(items, item)
		}

		counter++
	}

	return items, nil
}

func New(v app.YouTubeFeedEntry, dependencies YTDLDependencies) (Item, error) {

	item := Item{
		ID:          v.ID,
		GUID:        v.URL.String(),
		Title:       v.Title,
		Link:        &v.URL,
		Description: v.Description,
		PubDate:     v.Published,
		FileType:    "mp3",
		IsExplicit:  false,
	}

	err := item.enrich(dependencies.GetFileInformation)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (item *Item) enrich(fileGetter fileInformationGetter) error {
	details, err := fileGetter(item.ID)
	if err != nil {
		return err
	}

	item.FileURL = details.FileUrl
	item.Duration = details.Duration
	item.FileLength = details.ContentLength

	return nil
}

func (item Item) IsValid() bool {
	if item.Title == "" || item.Description == "" {
		return false
	}

	return true
}

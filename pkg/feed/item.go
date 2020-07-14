package feed

import (
	"github.com/psmarcin/youtubegoespodcast/pkg/video"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
	"github.com/rylio/ytdl"
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

type fileInformationGetter func(string, video.FileInformationGetter, video.FileUrlGetter) (video.Video, error)

func NewMap(videos []youtube.Video) ([]Item, error) {
	stream := make(chan Item, len(videos))
	for _, v := range videos {
		go func(video youtube.Video) {
			item, err := New(video)
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

func New(v youtube.Video) (Item, error) {
	u, err := url.Parse(v.Url)
	if err != nil {
		return Item{}, err
	}

	item := Item{
		ID:          v.ID,
		GUID:        v.Url,
		Title:       v.Title,
		Link:        u,
		Description: v.Description,
		PubDate:     v.PublishedAt,
		FileType:    "mp3",
		IsExplicit:  false,
	}

	err = item.enrich(video.GetFileInformation)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (item *Item) enrich(fileGetter fileInformationGetter) error {
	details, err := fileGetter(item.ID, ytdl.DefaultClient, ytdl.DefaultClient)
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

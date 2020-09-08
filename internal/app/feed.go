package app

import (
	"fmt"
	"github.com/eduncan911/podcast"
	feed2 "github.com/psmarcin/youtubegoespodcast/internal/domain/feed"
	"github.com/rylio/ytdl"
	"net/url"
	"os"
	"time"
)

const (
	Generator = "YouTubeGoesPodcast/v2"
)

type YouTubeDependency interface {
	ListEntry(string) ([]YouTubeFeedEntry, error)
	GetChannelCache(string) (YouTubeChannel, error)
}

type YTDLDependencies interface {
	GetFileURL(info *ytdl.VideoInfo) (url.URL, error)
	GetFileInformation(videoId string) (YTDLVideo, error)
}

type FeedService struct {
	youtubeService YouTubeDependency
	ytdlService    YTDLDependencies
}

type Feed struct {
	ChannelID string
	Content   podcast.Podcast
}

type FeedItem struct {
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

func NewFeedService(
	youtubeService YouTubeDependency,
	ytdlService YTDLDependencies,
) FeedService {
	return FeedService{
		youtubeService: youtubeService,
		ytdlService:    ytdlService,
	}
}

func (f FeedService) Create(channelID string) (Feed, error) {
	podcastFeed, err := f.GetDetails(channelID)
	feed := Feed{
		ChannelID: channelID,
		Content:   podcastFeed,
	}
	if err != nil {
		l.WithError(err).Errorf("can't get feed details for %s", channelID)
		return feed, err
	}

	videos, err := f.youtubeService.ListEntry(feed.ChannelID)
	if err != nil {
		l.WithError(err).Errorf("can't get video list for %s", channelID)
		return feed, err
	}

	items, err := f.CreateItems(videos)
	if err != nil {
		l.WithError(err).Errorf("can't enrich videos")
		return feed, err
	}

	lastPubDate := getLatestPubDate(items).Format(time.RFC1123Z)
	feed.Content.PubDate = lastPubDate
	feed.Content.LastBuildDate = lastPubDate

	feed.Content.Items, err = mapFeedItemToPodcastItem(items)

	if err != nil {
		l.WithError(err).Errorf("can't set videos for %s", channelID)
		return feed, err
	}

	feed.Content.Items = feed2.SortByPubDate(feed.Content.Items)
	feed.Content.Generator = Generator

	return feed, nil
}

func (f *FeedService) GetDetails(channelID string) (podcast.Podcast, error) {
	var fee podcast.Podcast
	channel, err := f.youtubeService.GetChannelCache(channelID)
	if err != nil {
		l.WithError(err).Error("can't get channel details")
		return fee, err
	}

	fee = podcast.New(channel.Title, channel.Url, channel.Description, &channel.PublishedAt, &channel.PublishedAt)

	// TODO: map categories from youtube to apple categories
	fee.AddCategory("Arts", []string{"Design"})
	fee.Language = channel.Country
	fee.Image = &podcast.Image{
		URL:         channel.Thumbnail,
		Title:       channel.Title,
		Link:        channel.Url,
		Description: channel.Description,
		//TODO: add width and height for thumbnail
		Width:  0,
		Height: 0,
	}
	fee.AddAuthor(channel.Author, channel.AuthorEmail)

	return fee, nil
}
func (f FeedService) CreateItems(items []YouTubeFeedEntry) ([]FeedItem, error) {
	stream := make(chan FeedItem, len(items))
	for _, v := range items {
		go func(video YouTubeFeedEntry) {
			item := FeedItem{
				ID:          v.ID,
				GUID:        v.URL.String(),
				Title:       v.Title,
				Link:        &v.URL,
				Description: v.Description,
				PubDate:     v.Published,
				FileType:    "mp3",
				IsExplicit:  false,
			}

			item, err := f.Enrich(item)
			if err != nil {
				l.WithError(err).WithField("video", video).Infof("can't create new item from video")
			}
			stream <- item
		}(v)
	}

	var is []FeedItem
	counter := 0
	for {
		if counter >= len(items) {
			break
		}

		item := <-stream
		if item.IsValid() {
			is = append(is, item)
		}

		counter++
	}

	return is, nil
}

func (f *FeedService) Enrich(item FeedItem) (FeedItem, error) {
	details, err := f.ytdlService.GetFileInformation(item.ID)
	if err != nil {
		return item, err
	}

	item.FileURL = details.FileUrl
	item.Duration = details.Duration
	item.FileLength = details.ContentLength

	return item, nil
}

func (f *Feed) Serialize() ([]byte, error) {
	return f.Content.Bytes(), nil
}

func getLatestPubDate(sourceItems []FeedItem) time.Time {
	// set channel last updated at field as latest item publishing date
	if len(sourceItems) != 0 {
		lastItem := sourceItems[0]
		return lastItem.PubDate
	}
	return time.Now()
}

func mapFeedItemToPodcastItem(sourceItems []FeedItem) ([]*podcast.Item, error) {
	var items []*podcast.Item

	for _, item := range sourceItems {
		videoURL := os.Getenv("API_URL") + "video/" + item.ID + "/track.mp3"

		items = append(items, &podcast.Item{
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

	}

	return items, nil
}

func (item FeedItem) IsValid() bool {
	if item.Title == "" || item.Description == "" {
		return false
	}

	return true
}

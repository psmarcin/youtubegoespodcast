package app

import (
	"context"
	"fmt"
	"github.com/eduncan911/podcast"
	feedDomain "github.com/psmarcin/youtubegoespodcast/internal/domain/feed"
	"github.com/rylio/ytdl"
	"go.opentelemetry.io/otel/label"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	Generator = "YouTubeGoesPodcast/v2"
)

type YouTubeDependency interface {
	ListEntry(context.Context, string) ([]YouTubeFeedEntry, error)
	GetChannelCache(context.Context, string) (YouTubeChannel, error)
}

type YTDLDependencies interface {
	GetFileURL(ctx context.Context, info *ytdl.VideoInfo) (url.URL, error)
	GetFileInformation(ctx context.Context, videoId string) (YTDLVideo, error)
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

func (f FeedService) Create(ctx context.Context, channelID string) (Feed, error) {
	ctx, span := tracer.Start(ctx, "feed-create")
	span.SetAttributes(label.Any("channelID", channelID))
	defer span.End()

	podcastFeed, err := f.GetFeedInformation(ctx, channelID)
	feed := Feed{
		ChannelID: channelID,
		Content:   podcastFeed,
	}
	if err != nil {
		l.WithError(err).Errorf("can't get feed details for %s", channelID)
		span.RecordError(ctx, err)
		return feed, err
	}

	videos, err := f.youtubeService.ListEntry(ctx, feed.ChannelID)
	if err != nil {
		l.WithError(err).Errorf("can't get video list for %s", channelID)
		span.RecordError(ctx, err)
		return feed, err
	}

	items, err := f.CreateItems(ctx, videos)
	if err != nil {
		l.WithError(err).Errorf("can't enrich videos")
		span.RecordError(ctx, err)
		return feed, err
	}

	lastPubDate := getLatestPubDate(items).Format(time.RFC1123Z)
	feed.Content.PubDate = lastPubDate
	feed.Content.LastBuildDate = lastPubDate

	feed.Content.Items, err = mapFeedItemToPodcastItem(items)

	if err != nil {
		l.WithError(err).Errorf("can't set videos for %s", channelID)
		span.RecordError(ctx, err)
		return feed, err
	}

	feed.Content.Items = feedDomain.SortByPubDate(feed.Content.Items)
	feed.Content.Generator = Generator

	return feed, nil
}

func (f *FeedService) GetFeedInformation(ctx context.Context, channelID string) (podcast.Podcast, error) {
	ctx, span := tracer.Start(ctx, "feed-get-information")
	span.SetAttributes(label.Any("channelID", channelID))
	defer span.End()

	var fee podcast.Podcast
	channel, err := f.youtubeService.GetChannelCache(ctx, channelID)
	if err != nil {
		l.WithError(err).Error("can't get channel details")
		span.RecordError(ctx, err)
		return fee, err
	}

	fee = podcast.New(channel.Title, channel.Url, channel.Description, &channel.PublishedAt, &channel.PublishedAt)

	lang := strings.ToLower(channel.Country)
	if lang == "" {
		lang = "en-us"
	}
	fee.Language = lang
	fee.AddCategory(channel.Category, []string{})
	fee.AddImage(channel.Thumbnail.Url.String())
	fee.Image = &podcast.Image{
		URL:         channel.Thumbnail.Url.String(),
		Title:       channel.Title,
		Link:        channel.Url,
		Description: channel.Description,
		Width:       channel.Thumbnail.Width,
		Height:      channel.Thumbnail.Height,
	}
	fee.AddAuthor(channel.Author, channel.AuthorEmail)
	fee.IExplicit = "false"

	return fee, nil
}
func (f FeedService) CreateItems(ctx context.Context, items []YouTubeFeedEntry) ([]FeedItem, error) {
	ctx, span := tracer.Start(ctx, "feed-create-items")
	span.SetAttributes(label.Any("items", label.ArrayValue(items)))
	defer span.End()

	stream := make(chan FeedItem, len(items))
	for _, v := range items {
		go func(video YouTubeFeedEntry) {
			item := FeedItem{
				ID:          video.ID,
				GUID:        video.URL.String(),
				Title:       video.Title,
				Link:        &video.URL,
				Description: video.Description,
				PubDate:     video.Published,
				FileType:    "mp3",
				IsExplicit:  false,
			}

			item, err := f.Enrich(ctx, item)
			if err != nil {
				l.WithError(err).WithField("video", video).Infof("can't create new item from video")
				span.RecordError(ctx, err)
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

func (f *FeedService) Enrich(ctx context.Context, item FeedItem) (FeedItem, error) {
	ctx, span := tracer.Start(ctx, "feed-enrich")
	span.SetAttributes(label.Any("item", item.ID))
	defer span.End()

	details, err := f.ytdlService.GetFileInformation(ctx, item.ID)
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

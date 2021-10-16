package app

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/eduncan911/podcast"
	feedDomain "github.com/psmarcin/youtubegoespodcast/internal/domain/feed"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	Generator = "YouTubeGoesPodcast/v2"
)

type YouTubeDependency interface {
	ListEntry(context.Context, string) ([]YouTubeFeedEntry, error)
	GetChannelCache(context.Context, string) (YouTubeChannel, error)
}

type YTDLDependencies interface {
	GetDetails(ctx context.Context, videoID string) (Details, error)
}

type FeedService struct {
	youtubeService YouTubeDependency
	fileService    YTDLDependencies
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
	fileService YTDLDependencies,
) FeedService {
	return FeedService{
		youtubeService: youtubeService,
		fileService:    fileService,
	}
}

func (f FeedService) Create(ctx context.Context, channelID string) (Feed, error) {
	ctx, span := tracer.Start(ctx, "feed-create")
	span.SetAttributes(attribute.String("channelID", channelID))
	defer span.End()

	podcastFeed, err := f.GetFeedInformation(ctx, channelID)
	feed := Feed{
		ChannelID: channelID,
		Content:   podcastFeed,
	}
	if err != nil {
		l.WithError(err).Errorf("can't get feed details for %s", channelID)
		span.RecordError(err)
		return feed, err
	}

	videos, err := f.youtubeService.ListEntry(ctx, feed.ChannelID)
	if err != nil {
		l.WithError(err).Errorf("can't get video list for %s", channelID)
		span.RecordError(err)
		return feed, err
	}

	items, err := f.CreateItems(ctx, videos)
	if err != nil {
		l.WithError(err).Errorf("can't enrich videos")
		span.RecordError(err)
		return feed, err
	}

	lastPubDate := getLatestPubDate(items).Format(time.RFC1123Z)
	feed.Content.PubDate = lastPubDate
	feed.Content.LastBuildDate = lastPubDate

	feed.Content.Items = mapFeedItemToPodcastItem(items)
	feed.Content.Items = feedDomain.SortByPubDate(feed.Content.Items)
	feed.Content.Generator = Generator

	return feed, nil
}

func (f *FeedService) GetFeedInformation(ctx context.Context, channelID string) (podcast.Podcast, error) {
	ctx, span := tracer.Start(ctx, "feed-get-information")
	span.SetAttributes(attribute.String("channelID", channelID))
	defer span.End()

	var fee podcast.Podcast
	channel, err := f.youtubeService.GetChannelCache(ctx, channelID)
	if err != nil {
		l.WithError(err).Error("can't get channel details")
		span.RecordError(err)
		return fee, err
	}

	fee = podcast.New(channel.Title, channel.URL, channel.Description, &channel.PublishedAt, &channel.PublishedAt)

	lang := strings.ToLower(channel.Country)
	if lang == "" {
		lang = "en-us"
	}
	fee.Language = lang
	fee.AddCategory(channel.Category, []string{})
	fee.AddImage(channel.Thumbnail.URL.String())
	fee.Image = &podcast.Image{
		URL:         channel.Thumbnail.URL.String(),
		Title:       channel.Title,
		Link:        channel.URL,
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
	span.SetAttributes(attribute.Int("itemsCount", len(items)))
	defer span.End()

	stream := make(chan FeedItem, len(items))
	for _, v := range items {
		go func(video YouTubeFeedEntry) {
			span.AddEvent("add-item", trace.WithAttributes(attribute.KeyValue{
				Key:   attribute.Key(fmt.Sprintf("video-%s", video.ID)),
				Value: attribute.StringValue(video.URL.String()),
			}))
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
	span.SetAttributes(attribute.String("item", item.ID))
	defer span.End()

	details, err := f.fileService.GetDetails(ctx, item.ID)
	if err != nil {
		span.RecordError(err)
		return item, err
	}

	item.FileURL = &details.URL
	item.Duration = 0
	item.FileLength = 0

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

func mapFeedItemToPodcastItem(sourceItems []FeedItem) []*podcast.Item {
	var items []*podcast.Item

	for _, item := range sourceItems {
		videoURL := os.Getenv("API_URL") + "video/" + item.ID + "/track.mp3"

		pubDate := item.PubDate

		items = append(items, &podcast.Item{
			GUID:        item.GUID,
			Title:       item.Title,
			Link:        item.Link.String(),
			Description: item.Description,
			PubDate:     &pubDate,
			Enclosure: &podcast.Enclosure{
				URL:    videoURL,
				Length: item.FileLength,
				Type:   podcast.MP3,
			},
			IDuration: fmt.Sprintf("%f", item.Duration.Seconds()),
			IExplicit: "no",
		})
	}

	return items
}

func (item FeedItem) IsValid() bool {
	if item.Title == "" || item.Description == "" {
		return false
	}

	return true
}

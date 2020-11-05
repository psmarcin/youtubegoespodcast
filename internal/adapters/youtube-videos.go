package adapters

import (
	"context"
	"encoding/xml"
	"github.com/pkg/errors"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"go.opentelemetry.io/otel/label"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	BaseUrl     = "https://www.youtube.com/"
	FeedUrlBase = BaseUrl + "feeds/videos.xml"
)

var FeedUrl *url.URL

type Feed struct {
	Entry []FeedEntry `xml:"entry"`
}

type FeedEntry struct {
	ID      string `xml:"id"`
	VideoId string `xml:"videoId"`
	Title   string `xml:"title"`
	Link    struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Published string `xml:"published"`
	Group     struct {
		Thumbnail struct {
			URL    string `xml:"url,attr"`
			Width  string `xml:"width,attr"`
			Height string `xml:"height,attr"`
		} `xml:"thumbnail"`
		Description string `xml:"description"`
	} `xml:"group"`
}

type YouTubeRepository struct {
	FeedURL *url.URL
}

func NewYouTubeRepository() (YouTubeRepository, error) {
	feedUrl, err := url.Parse(FeedUrlBase)
	if err != nil {
		return YouTubeRepository{}, errors.Wrap(err, "can't parse feed base url")
	}
	return YouTubeRepository{
		FeedURL: feedUrl,
	}, nil
}
func init() {
	feedUrl, err := url.Parse(FeedUrlBase)
	if err != nil {
		l.WithError(err).Fatalf("can't parse feed base url")
	}
	FeedUrl = feedUrl
}

func (y YouTubeRepository) ListEntry(ctx context.Context, channelId string) ([]app.YouTubeFeedEntry, error) {
	var entries []app.YouTubeFeedEntry
	var f Feed
	ctx, span := tracer.Start(ctx, "list-entry")
	span.SetAttributes(label.String("channel-id", channelId))
	defer span.End()

	u := *FeedUrl
	q := u.Query()
	q.Add("channel_id", channelId)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		l.WithError(err).Errorf("can't get channel feed")
		span.RecordError(ctx, err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Errorf("can't parse channel feed")
		span.RecordError(ctx, err)
		return nil, err
	}

	err = xml.Unmarshal(body, &f)
	if err != nil {
		l.WithError(err).WithField("body", string(body)).Errorf("can't unmarshal channel feed")
		span.RecordError(ctx, err)
		return nil, err
	}

	for _, entry := range f.Entry {
		ytFeedEntry, err := mapFeedToYouTubeEntry(entry)
		if err != nil {
			span.RecordError(ctx, err)
			return entries, errors.Wrapf(err, "unable to map FeedEntry to YouTubeEntry")
		}
		entries = append(entries, ytFeedEntry)
	}

	return entries, nil

}

func mapFeedToYouTubeEntry(entry FeedEntry) (app.YouTubeFeedEntry, error) {
	ytFeedEntry := app.YouTubeFeedEntry{}
	publishedAt, err := time.Parse(time.RFC3339, entry.Published)
	if err != nil {
		return ytFeedEntry, errors.Wrapf(err, "unable to parse publish date: %s for video id: %s", entry.Published, entry.ID)
	}

	u, err := url.Parse(entry.Link.Href)
	if err != nil {
		return ytFeedEntry, errors.Wrapf(err, "unable to parse url: %s for video id: %s", entry.Link.Href, entry.ID)
	}

	tUrl, err := url.Parse(entry.Link.Href)
	if err != nil {
		l.WithError(err).Errorf("can't parse video %s thumbnail url %s", entry.ID, entry.Group.Thumbnail.URL)
	}

	tHeight, err := strconv.Atoi(entry.Group.Thumbnail.Height)
	if err != nil {
		l.WithError(err).Errorf("can't parse video %s thumbnail height %s", entry.ID, entry.Group.Thumbnail.Height)
	}

	tWidth, err := strconv.Atoi(entry.Group.Thumbnail.Width)
	if err != nil {
		l.WithError(err).Errorf("can't parse video %s thumbnail width %s", entry.ID, entry.Group.Thumbnail.Width)
	}

	ytFeedEntry = app.YouTubeFeedEntry{
		ID:          entry.VideoId,
		Title:       entry.Title,
		Description: entry.Group.Description,
		Published:   publishedAt,
		URL:         *u,
		Thumbnail: app.YouTubeThumbnail{
			Height: tHeight,
			Width:  tWidth,
			Url:    *tUrl,
		},
	}

	return ytFeedEntry, nil
}

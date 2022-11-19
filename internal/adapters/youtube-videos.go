package adapters

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/psmarcin/youtubegoespodcast/internal/app"
	"go.opentelemetry.io/otel/attribute"
)

const (
	BaseURL           = "https://www.youtube.com/"
	FeedURLBase       = BaseURL + "feeds/videos.xml"
	HTTPClientTimeout = time.Second * 5
)

var FeedURL *url.URL

type Feed struct {
	Entry []FeedEntry `xml:"entry"`
}

type FeedEntry struct {
	ID      string `xml:"id"`
	VideoID string `xml:"videoId"`
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
	FeedURL    *url.URL
	HTTPClient *http.Client
}

func NewYouTubeRepository() (YouTubeRepository, error) {
	feedURL, err := url.Parse(FeedURLBase)
	FeedURL = feedURL
	if err != nil {
		return YouTubeRepository{}, errors.Wrap(err, "can't parse feed base url")
	}

	client := &http.Client{
		Timeout: HTTPClientTimeout,
	}

	return YouTubeRepository{
		FeedURL:    feedURL,
		HTTPClient: client,
	}, nil
}

func (y YouTubeRepository) ListEntry(ctx context.Context, channelID string) ([]app.YouTubeFeedEntry, error) {
	var entries []app.YouTubeFeedEntry
	var f Feed
	_, span := tracer.Start(ctx, "list-entry")
	span.SetAttributes(attribute.String("channel-id", channelID))
	defer span.End()

	u := *FeedURL
	q := u.Query()
	q.Add("channel_id", channelID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		l.WithError(err).Errorf("can't get create request to get feed")
		span.RecordError(err)
		return nil, err
	}

	resp, err := y.HTTPClient.Do(req)
	if err != nil {
		l.WithError(err).Errorf("can't do request to get feed")
		span.RecordError(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Errorf("can't parse channel feed")
		span.RecordError(err)
		return nil, err
	}

	err = xml.Unmarshal(body, &f)
	if err != nil {
		l.WithError(err).WithField("body", string(body)).Errorf("can't unmarshal channel feed")
		span.RecordError(err)
		return nil, err
	}

	for _, entry := range f.Entry {
		ytFeedEntry, err := mapFeedToYouTubeEntry(entry)
		if err != nil {
			span.RecordError(err)
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

	tURL, err := url.Parse(entry.Link.Href)
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
		ID:          entry.VideoID,
		Title:       entry.Title,
		Description: entry.Group.Description,
		Published:   publishedAt,
		URL:         *u,
		Thumbnail: app.YouTubeThumbnail{
			Height: tHeight,
			Width:  tWidth,
			URL:    *tURL,
		},
	}

	return ytFeedEntry, nil
}

package youtube

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	BASE_URL      = "https://www.youtube.com/"
	FEED_URL_BASE = BASE_URL + "feeds/videos.xml"
)

var FEED_URL *url.URL

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Text    string   `xml:",chardata"`
	Yt      string   `xml:"yt,attr"`
	Media   string   `xml:"media,attr"`
	Xmlns   string   `xml:"xmlns,attr"`
	Link    []struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	ID        string `xml:"id"`
	ChannelId string `xml:"channelId"`
	Title     string `xml:"title"`
	Author    struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
		URI  string `xml:"uri"`
	} `xml:"author"`
	Published string `xml:"published"`
	Entry     []struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id"`
		VideoId   string `xml:"videoId"`
		ChannelId string `xml:"channelId"`
		Title     string `xml:"title"`
		Link      struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Author struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
			URI  string `xml:"uri"`
		} `xml:"author"`
		Published string `xml:"published"`
		Updated   string `xml:"updated"`
		Group     struct {
			Text    string `xml:",chardata"`
			Title   string `xml:"title"`
			Content struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Type   string `xml:"type,attr"`
				Width  string `xml:"width,attr"`
				Height string `xml:"height,attr"`
			} `xml:"content"`
			Thumbnail struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Width  string `xml:"width,attr"`
				Height string `xml:"height,attr"`
			} `xml:"thumbnail"`
			Description string `xml:"description"`
			Community   struct {
				Text       string `xml:",chardata"`
				StarRating struct {
					Text    string `xml:",chardata"`
					Count   string `xml:"count,attr"`
					Average string `xml:"average,attr"`
					Min     string `xml:"min,attr"`
					Max     string `xml:"max,attr"`
				} `xml:"starRating"`
				Statistics struct {
					Text  string `xml:",chardata"`
					Views string `xml:"views,attr"`
				} `xml:"statistics"`
			} `xml:"community"`
		} `xml:"group"`
	} `xml:"entry"`
}

func init() {
	feedUrl, err := url.Parse(FEED_URL_BASE)
	if err != nil {
		l.WithError(err).Fatalf("can't parse feed base url")
	}
	FEED_URL = feedUrl
}

func GetFeed(channelId string) (Feed, error) {
	var f Feed
	u := *FEED_URL
	q := u.Query()
	q.Add("channel_id", channelId)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		l.WithError(err).Errorf("can't get channel feed")
		return f, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.WithError(err).Errorf("can't parse channel feed")
		return f, err
	}

	err = xml.Unmarshal(body, &f)
	if err != nil {
		l.WithError(err).WithField("body", body).Errorf("can't unmarshal channel feed")
		return f, err
	}

	return f, nil

}

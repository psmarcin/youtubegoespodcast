package feed

import (
	"encoding/xml"
)

const (
	ytChannelURL = "https://youtube.com/channel/"
)

type Feed struct {
	XMLName       string     `xml:"channel"`
	ChannelID     string     `xml:"-"`
	Title         string     `xml:"title"`
	Link          string     `xml:"link"`
	Description   string     `xml:"description"`
	Category      string     `xml:"category"`
	Generator     string     `xml:"generator"`
	Language      string     `xml:"language"`
	LastBuildDate string     `xml:"lastBuildDate"`
	PubDate       string     `xml:"pubDate"`
	Image         Image      `xml:"image"`
	ITAuthor      string     `xml:"itunes:author"`
	ITSubtitle    string     `xml:"itunes:subtitle"`
	ITSummary     ITSummary  `xml:"itunes:summary"`
	ITImage       ITImage    `xml:"itunes:image"`
	ITExplicit    string     `xml:"itunes:explicit"`
	ITCategory    ITCategory `xml:"itunes:category"`
	Items         []Item     `xml:"item"`
}
type Image struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}
type ITImage struct {
	Href string `xml:"href,attr"`
}
type ITCategory struct {
	Text string `xml:"text,attr"`
}
type ITSummary struct {
	XMLName xml.Name `xml:"itunes:summary"`
	Text    string   `xml:",cdata"`
}

type Item struct {
	GUID        string    `xml:"guid"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     string    `xml:"pubDate"`
	Enclosure   Enclosure `xml:"enclosure"`
	ITAuthor    string    `xml:"itunes:author"`
	ITSubtitle  string    `xml:"itunes:subtitle"`
	ITSummary   ITSummary `xml:"itunes:summary"`
	ITImage     ITImage   `xml:"itunes:image"`
	ITDuration  string    `xml:"itunes:duration"`
	ITExplicit  string    `xml:"itunes:explicit"`
	ITOrder     string    `xml:"itunes:order"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

func (f *Feed) serialize() ([]byte, error) {
	xml, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return []byte{}, err
	}
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
			<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
			` + string(xml) + `
			</rss>`), nil
}

func new(channelID string) Feed {
	feed := Feed{
		XMLName:   "channel",
		ChannelID: channelID,
		Generator: "ygp v1",
	}
	return feed
}

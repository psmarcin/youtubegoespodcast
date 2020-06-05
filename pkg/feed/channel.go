package feed

import (
	"github.com/eduncan911/podcast"
	"github.com/psmarcin/youtubegoespodcast/pkg/youtube"
)

func (f *Feed) AddItem(item podcast.Item) error {
	if item.Title != "" && item.Enclosure.URL != "" {
		_, err := f.Content.AddItem(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Feed) GetDetails(channelID string) error {
	channel, err := youtube.Yt.ChannelsGetFromCache(channelID)
	if err != nil {
		l.WithError(err).Error("can't get channel details")
		return err
	}

	fee := podcast.New(channel.Title, channel.Url, channel.Description, &channel.PublishedAt, &channel.PublishedAt)

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
	f.Content = fee

	return nil
}

package feed

import (
	"errors"
	"fmt"
	"github.com/eduncan911/podcast"
	"net/http"
	"os"
	"strconv"
	"time"

	vid "ygp/pkg/video"
	"ygp/pkg/youtube"
)

type VideoFileDetails struct {
	ContentType   string
	ContentLength string
}

func (f *Feed) SetVideos(videos []youtube.Video) error {
	filteredVideos := filterEmptyVideoOut(videos)
	// set channel last updated at field as latest item publishing date
	if len(filteredVideos) != 0 {
		lastItem := videos[0]
		f.Content.LastBuildDate = lastItem.PublishedAt.Format(time.RFC1123Z)
		f.Content.PubDate = f.Content.LastBuildDate
	}

	stream := make(chan podcast.Item, len(filteredVideos))
	for i, video := range videos {
		go func(video youtube.Video, i int) error {
			vd, err := vid.GetDetails(video.ID)
			videoURL := os.Getenv("API_URL") + "video/" + video.ID + "/track.mp3"
			fileDetails, _ := GetVideoFileDetails(videoURL)

			if err != nil {
				l.WithError(err).Error("item")
				stream <- podcast.Item{}
				return err
			}

			encodeLength, err := strconv.ParseInt(fileDetails.ContentLength, 10, 32)
			if err != nil {
				stream <- podcast.Item{}
				return err
			}

			it := podcast.Item{
				GUID:        video.ID,
				Title:       video.Title,
				Link:        video.Url,
				Description: video.Description,
				PubDate:     &video.PublishedAt,
				Enclosure: &podcast.Enclosure{
					URL:    videoURL,
					Length: encodeLength,
					Type:   podcast.MP3,
				},
				IDuration: fmt.Sprintf("%f", vd.Duration.Seconds()),
				IExplicit: "no",
				//IOrder:    fmt.Sprintf("%d", i+1),
			}
			stream <- it
			return nil
		}(video, i)

	}

	counter := 0
	for {
		if counter >= len(filteredVideos) {
			break
		}
		err := f.AddItem(<-stream)
		if err != nil {
			return err
		}
		counter++
	}

	return nil
}

func GetVideoFileDetails(videoURL string) (VideoFileDetails, error) {
	resp, err := http.Head(videoURL)
	if err != nil {
		return VideoFileDetails{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return VideoFileDetails{}, errors.New("[ITEM] Can't get file details for " + videoURL)
	}
	return VideoFileDetails{
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.Header.Get("Content-Length"),
	}, nil
}

func filterEmptyVideoOut(videos []youtube.Video) []youtube.Video{
	var filtered []youtube.Video
	for _, video := range videos{
		if video.ID == ""{
			continue
		}
		filtered = append(filtered, video)
	}

	return filtered
}

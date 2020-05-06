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

	"github.com/sirupsen/logrus"
)

const (
	ytVideoURL = "https://youtube.com/watch?v="
)

type VideosResponse struct {
	RegionCode string        `json:"regionCode"`
	Items      []VideosItems `json:"items"`
}
type VideosItems struct {
	ID      ID            `json:"id"`
	Snippet VideosSnippet `json:"snippet"`
}
type ID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}
type VideosSnippet struct {
	PublishedAt          time.Time                `json:"publishedAt"`
	ChannelID            string                   `json:"channelId"`
	Title                string                   `json:"title"`
	Description          string                   `json:"description"`
	Thumbnails           ChannelDetailsThumbnails `json:"thumbnails"`
	ChannelTitle         string                   `json:"channelTitle"`
	LiveBroadcastContent string                   `json:"liveBroadcastContent"`
}

type VideoFileDetails struct {
	ContentType   string
	ContentLength string
}

type VideoDetails struct {
	Duration time.Duration
}

type VideosDetailsResponse struct {
	Items []VideosDetailsItems `json:"items"`
}
type VideosDetailsItems struct {
	Details VideosDetailsContent `json:"contentDetails"`
}
type VideosDetailsContent struct {
	Duration string `json:"duration"`
}

func (f *Feed) GetVideos(q string) (VideosResponse, error) {
	videos := VideosResponse{}

	req, err := http.NewRequest("GET", youtube.YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideosResponse{}, err
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("order", "date")
	query.Add("channelId", f.ChannelID)
	query.Add("maxResults", "25")
	query.Add("q", q)
	query.Add("fields", "items(id,snippet(channelId,channelTitle,description,publishedAt,thumbnails/high,title))")
	req.URL.RawQuery = query.Encode()

	requestErr := youtube.Request(req, &videos)

	return videos, requestErr
}

func (f *Feed) SetVideos(videos VideosResponse) error {
	stream := make(chan podcast.Item, countItems(videos.Items))

	// set channel last updated at field as latest item publishing date
	if len(videos.Items) != 0 {
		lastItem := videos.Items[0]
		f.Content.LastBuildDate = lastItem.Snippet.PublishedAt.Format(time.RFC1123Z)
		f.Content.PubDate = f.Content.LastBuildDate
	}

	for i, video := range videos.Items {
		s := video.Snippet

		go func(video VideosItems, i int) error {

			if video.ID.VideoID == "" {
				stream <- podcast.Item{}
				return errors.New("[ITEM] video id not provided for set video")
			}

			vd, err := vid.GetDetails(video.ID.VideoID)
			videoURL := os.Getenv("API_URL") + "video/" + video.ID.VideoID + "/track.mp3"
			fileDetails, _ := GetVideoFileDetails(videoURL)

			if err != nil {
				logrus.WithError(err).Printf("[ITEM]")
				stream <- podcast.Item{}
				return err
			}

			encodeLength, erro := strconv.ParseInt(fileDetails.ContentLength, 10, 32)

			if erro != nil {
				stream <- podcast.Item{}
				return erro
			}
			it := podcast.Item{
				GUID:        video.ID.VideoID,
				Title:       s.Title,
				Link:        ytVideoURL + video.ID.VideoID,
				Description: s.Description,
				PubDate:     &s.PublishedAt,
				Enclosure: &podcast.Enclosure{
					URL:    videoURL,
					Length: encodeLength,
					Type:   podcast.MP3,
				},
				IDuration: fmt.Sprintf("%f", vd.Duration.Seconds()),
				IExplicit: "no",
				IOrder:    fmt.Sprintf("%d", i+1),
			}
			stream <- it
			return nil
		}(video, i)

	}

	counter := 0
	for {
		if counter >= len(videos.Items) {
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

func countItems(items []VideosItems) int {
	counter := 0
	for _, item := range items {
		if item.ID.VideoID != "" {
			counter++
		}
	}
	return counter
}

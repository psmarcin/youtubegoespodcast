package feed

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"ytg/pkg/redis_client"
	"ytg/pkg/youtube"

	"github.com/sirupsen/logrus"
)

const (
	ytVideoURL            = "https://youtube.com/watch?v="
	videoItemCachePrefix  = "ytvideo_"
	videoItemCacheTTL     = time.Hour * 24 * 7
	videoItemsCachePrefix = "ytvideos_"
	videoItemsCacheTTL    = time.Hour
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

func (f *Feed) getVideos(q string) (VideosResponse, error) {
	videos := VideosResponse{}

	_, err := redis_client.Client.GetKey(videoItemsCachePrefix+f.ChannelID, &videos)
	// got cached value, fast return
	if err == nil {
		return videos, nil
	}

	req, err := http.NewRequest("GET", youtube.YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideosResponse{}, err
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("order", "date")
	query.Add("channelId", f.ChannelID)
	query.Add("maxResults", "10")
	query.Add("q", q)
	query.Add("fields", "items(id,snippet(channelId,channelTitle,description,publishedAt,thumbnails/high,title))")
	req.URL.RawQuery = query.Encode()

	err = youtube.Request(req, &videos)
	if err != nil {
		return VideosResponse{}, err
	}

	// save video items to cache
	str, err := json.Marshal(videos)
	if err != nil {
		return videos, err
	}
	go redis_client.Client.SetKey(videoItemsCachePrefix+f.ChannelID, string(str), videoItemsCacheTTL)
	return videos, nil
}

func (f *Feed) setVideos(videos VideosResponse) error {
	stream := make(chan Item, len(videos.Items))

	for i, video := range videos.Items {
		s := video.Snippet
		go func(video VideosItems, i int) error {
			vd, err := getVideoDetails(video.ID.VideoID)
			videoURL := os.Getenv("API_URL") + "video/" + video.ID.VideoID + ".mp4"
			fileDetails, _ := getVideoFileDetails(videoURL)
			if err != nil {
				logrus.WithError(err).Printf("[ITEM]")
				stream <- Item{}
				return err
			}

			stream <- Item{
				GUID:        video.ID.VideoID,
				Title:       s.Title,
				Link:        videoURL,
				Description: s.Description,
				PubDate:     s.PublishedAt.String(),
				Enclosure: Enclosure{
					URL:    videoURL,
					Length: fileDetails.ContentLength,
					Type:   fileDetails.ContentType,
				},
				ITAuthor:   f.ITAuthor,
				ITSubtitle: s.Title,
				ITSummary: ITSummary{
					Text: s.Description,
				},
				ITImage: ITImage{
					Href: getImageURL(s.Thumbnails.High.URL),
				},
				ITExplicit: "no",
				ITDuration: calculateDuration(vd.Duration),
			}
			return nil
		}(video, i)

	}
	counter := 0
	for {
		if counter >= len(videos.Items) {
			break
		}

		f.addItem(<-stream)
		counter++
	}
	return nil
}

func getVideoFileDetails(videoURL string) (VideoFileDetails, error) {
	resp, err := http.Head(videoURL)
	if err != nil {
		return VideoFileDetails{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return VideoFileDetails{}, errors.New("Can't get file details for " + videoURL)
	}
	return VideoFileDetails{
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.Header.Get("Content-Length"),
	}, nil
}

func getVideoDetails(videoID string) (VideoDetails, error) {
	vd := VideoDetails{}
	videoDetails := VideosDetailsResponse{}

	_, err := redis_client.Client.GetKey(videoItemCachePrefix+videoID, &vd)
	// got cached value, fast return
	if err == nil {
		return vd, nil
	}

	// no cached value, have to make call to API
	req, err := http.NewRequest("GET", youtube.YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideoDetails{}, err
	}
	query := req.URL.Query()
	query.Add("part", "contentDetails")
	query.Add("id", videoID)
	query.Add("maxResults", "1")
	query.Add("fields", "items/contentDetails/duration")
	req.URL.RawQuery = query.Encode()

	err = youtube.Request(req, &videoDetails)
	if err != nil {
		return VideoDetails{}, err
	}
	if len(videoDetails.Items) != 1 {
		return VideoDetails{}, errors.New("Can't get video details")
	}
	duration, err := parseDuration(videoDetails.Items[0].Details.Duration)
	if err != nil {
		return VideoDetails{}, err
	}
	vd.Duration = duration

	// save videoDetails to cache
	str, err := json.Marshal(vd)
	if err != nil {
		return vd, err
	}
	go redis_client.Client.SetKey(videoItemCachePrefix+videoID, string(str), videoItemCacheTTL)

	return vd, nil
}

func parseDuration(duration string) (time.Duration, error) {
	durationString := normalizeDurationString(duration)
	return time.ParseDuration(durationString)
}

func normalizeDurationString(duration string) string {
	return strings.ToLower(strings.Replace(duration, "PT", "", 1))
}

func getImageURL(src string) string {
	return src
}

func calculateDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

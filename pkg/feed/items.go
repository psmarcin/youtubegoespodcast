package feed

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"ygp/pkg/errx"
	"ygp/pkg/redis_client"
	"ygp/pkg/youtube"

	"github.com/sirupsen/logrus"
)

const (
	ytVideoURL            = "https://youtube.com/watch?v="
	videoItemCachePrefix  = "ytvideo_"
	videoItemCacheTTL     = time.Hour * 24 * 60
	videoItemsCachePrefix = "ytvideos_"
	videoItemsCacheTTL    = time.Hour * 24
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

func (f *Feed) getVideos(q string) (VideosResponse, errx.APIError) {
	videos := VideosResponse{}

	_, err := redis_client.Client.GetKey(videoItemsCachePrefix+f.ChannelID, &videos)
	// got cached value, fast return
	if err == nil {
		return videos, errx.APIError{}
	}

	req, err := http.NewRequest("GET", youtube.YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideosResponse{}, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	query := req.URL.Query()
	query.Add("part", "snippet")
	query.Add("order", "date")
	query.Add("channelId", f.ChannelID)
	query.Add("maxResults", "5")
	query.Add("q", q)
	query.Add("fields", "items(id,snippet(channelId,channelTitle,description,publishedAt,thumbnails/high,title))")
	req.URL.RawQuery = query.Encode()

	requestErr := youtube.Request(req, &videos)
	if requestErr.IsError() {
		return VideosResponse{}, requestErr
	}

	// save video items to cache
	str, err := json.Marshal(videos)
	if err != nil {
		return videos, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	go redis_client.Client.SetKey(videoItemsCachePrefix+f.ChannelID, string(str), videoItemsCacheTTL)
	return videos, errx.APIError{}
}

func (f *Feed) setVideos(videos VideosResponse) errx.APIError {
	stream := make(chan Item, countItems(videos.Items))

	for i, video := range videos.Items {
		s := video.Snippet

		go func(video VideosItems, i int) errx.APIError {

			if video.ID.VideoID == "" {
				stream <- Item{}
				return errx.APIError{}
			}

			vd, err := getVideoDetails(video.ID.VideoID)
			videoURL := os.Getenv("API_URL") + "video/" + video.ID.VideoID + ".mp4"
			fileDetails, _ := getVideoFileDetails(videoURL)
			if err.IsError() {
				logrus.WithError(err.Err).Printf("[ITEM]")
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
			return errx.APIError{}
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
	return errx.APIError{}
}

func getVideoFileDetails(videoURL string) (VideoFileDetails, errx.APIError) {
	resp, err := http.Head(videoURL)
	if err != nil {
		return VideoFileDetails{}, errx.NewAPIError(err, http.StatusBadRequest)
	}
	if resp.StatusCode != http.StatusOK {
		return VideoFileDetails{}, errx.NewAPIError(
			errors.New("Can't get file details for "+videoURL),
			http.StatusBadRequest,
		)
	}
	return VideoFileDetails{
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.Header.Get("Content-Length"),
	}, errx.APIError{}
}

func getVideoDetails(videoID string) (VideoDetails, errx.APIError) {
	vd := VideoDetails{}
	videoDetails := VideosDetailsResponse{}

	_, err := redis_client.Client.GetKey(videoItemCachePrefix+videoID, &vd)
	// got cached value, fast return
	if err == nil {
		return vd, errx.APIError{}
	}

	// no cached value, have to make call to API
	req, err := http.NewRequest("GET", youtube.YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideoDetails{}, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	query := req.URL.Query()
	query.Add("part", "contentDetails")
	query.Add("id", videoID)
	query.Add("maxResults", "1")
	query.Add("fields", "items/contentDetails/duration")
	req.URL.RawQuery = query.Encode()

	requestErr := youtube.Request(req, &videoDetails)
	if requestErr.IsError() {
		return VideoDetails{}, requestErr
	}
	if len(videoDetails.Items) != 1 {
		return VideoDetails{}, errx.NewAPIError(errors.New("Can't get video details for "+videoID), http.StatusBadRequest)
	}
	duration, err := parseDuration(videoDetails.Items[0].Details.Duration)
	if err != nil {
		return VideoDetails{}, errx.NewAPIError(errors.New("Can't parse duration for "+videoID), http.StatusInternalServerError)
	}
	vd.Duration = duration

	// save videoDetails to cache
	str, err := json.Marshal(vd)
	if err != nil {
		return vd, errx.NewAPIError(err, http.StatusInternalServerError)
	}
	go redis_client.Client.SetKey(videoItemCachePrefix+videoID, string(str), videoItemCacheTTL)

	return vd, errx.APIError{}
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

func countItems(items []VideosItems) int {
	counter := 0
	for _, item := range items {
		if item.ID.VideoID != "" {
			counter++
		}
	}
	return counter
}

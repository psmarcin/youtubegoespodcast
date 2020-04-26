package feed

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eduncan911/podcast"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"ygp/pkg/cache"
	"ygp/pkg/errx"
	"ygp/pkg/youtube"

	"github.com/sirupsen/logrus"
)

var _cacheClient = cache.Client

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

func (f *Feed) GetVideos(q string, disbleCache bool) (VideosResponse, errx.APIError) {
	videos := VideosResponse{}

	if !disbleCache {
		_, err := _cacheClient.GetKey(videoItemsCachePrefix+f.ChannelID, &videos)
		// got cached value, fast return
		if err == nil {
			return videos, errx.APIError{}
		}
	}

	req, err := http.NewRequest("GET", youtube.YouTubeURL+"search", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideosResponse{}, errx.New(err, http.StatusInternalServerError)
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

	if !disbleCache {
		// save video items to cache
		str, err := json.Marshal(videos)
		if err != nil {
			return videos, errx.New(err, http.StatusInternalServerError)
		}
		go _cacheClient.SetKey(videoItemsCachePrefix+f.ChannelID, string(str), videoItemsCacheTTL)

	}
	return videos, errx.APIError{}
}

func (f *Feed) SetVideos(videos VideosResponse) errx.APIError {
	stream := make(chan podcast.Item, countItems(videos.Items))

	// set channel last updated at field as latest item publishing date
	if len(videos.Items) != 0 {
		lastItem := videos.Items[0]
		f.Content.LastBuildDate = lastItem.Snippet.PublishedAt.Format(time.RFC1123Z)
		f.Content.PubDate = f.Content.LastBuildDate
	}

	for i, video := range videos.Items {
		s := video.Snippet

		go func(video VideosItems, i int) errx.APIError {

			if video.ID.VideoID == "" {
				stream <- podcast.Item{}
				return errx.APIError{}
			}

			vd, err := GetVideoDetails(video.ID.VideoID, false)
			videoURL := os.Getenv("API_URL") + "video/" + video.ID.VideoID + "/track.mp3"
			fileDetails, _ := GetVideoFileDetails(videoURL)
			if err.IsError() {
				logrus.WithError(err.Err).Printf("[ITEM]")
				stream <- podcast.Item{}
				return err
			}

			encodeLength, erro := strconv.ParseInt(fileDetails.ContentLength, 10, 32)

			if erro != nil {
				stream <- podcast.Item{}
				return errx.New(erro, 500)
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
			return errx.APIError{}
		}(video, i)

	}
	counter := 0
	for {
		if counter >= len(videos.Items) {
			break
		}
		err := f.AddItem(<-stream)
		if err != nil {
			return errx.New(err, 500)
		}
		counter++
	}
	return errx.APIError{}
}

func GetVideoFileDetails(videoURL string) (VideoFileDetails, errx.APIError) {
	resp, err := http.Head(videoURL)
	if err != nil {
		return VideoFileDetails{}, errx.New(err, http.StatusBadRequest)
	}
	if resp.StatusCode != http.StatusOK {
		return VideoFileDetails{}, errx.New(
			errors.New("Can't get file details for "+videoURL),
			http.StatusBadRequest,
		)
	}
	return VideoFileDetails{
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.Header.Get("Content-Length"),
	}, errx.APIError{}
}

func GetVideoDetails(videoID string, disableCache bool) (VideoDetails, errx.APIError) {
	vd := VideoDetails{}
	videoDetails := VideosDetailsResponse{}

	if !disableCache {
		_, err := _cacheClient.GetKey(videoItemCachePrefix+videoID, &vd)
		// got cached value, fast return
		if err == nil {
			return vd, errx.APIError{}
		}
	}
	// no cached value, have to make call to API
	req, err := http.NewRequest("GET", youtube.YouTubeURL+"videos", nil)
	if err != nil {
		logrus.WithError(err).Fatal("[YT] Can't create new request")
		return VideoDetails{}, errx.New(err, http.StatusInternalServerError)
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
		return VideoDetails{}, errx.New(errors.New("Can't get video details for "+videoID), http.StatusBadRequest)
	}
	duration, err := parseDuration(videoDetails.Items[0].Details.Duration)
	if err != nil {
		return VideoDetails{}, errx.New(errors.New("Can't parse duration for "+videoID), http.StatusInternalServerError)
	}
	vd.Duration = duration

	if !disableCache { // save videoDetails to cache
		str, err := json.Marshal(vd)
		if err != nil {
			return vd, errx.New(err, http.StatusInternalServerError)
		}
		go _cacheClient.SetKey(videoItemCachePrefix+videoID, string(str), videoItemCacheTTL)
	}

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

func countItems(items []VideosItems) int {
	counter := 0
	for _, item := range items {
		if item.ID.VideoID != "" {
			counter++
		}
	}
	return counter
}

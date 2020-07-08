package video

import (
	"errors"
	"github.com/psmarcin/youtubegoespodcast/pkg/cache"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	cacheKeyPrefix = "video_file_details_"
	cacheTTL       = time.Hour * 24 * 31
)

type FileDetails struct {
	ContentType   string `json:"contentType"`
	ContentLength int64  `json:"contentLength"`
}

func (f FileDetails) Get(url *url.URL) (FileDetails, error) {
	resp, err := http.Head(url.String())
	if err != nil {
		return f, err
	}
	if resp.StatusCode != http.StatusOK {
		return f, errors.New("[ITEM] Can't get file details for " + url.String())
	}
	ct := resp.Header.Get("Content-Type")
	cl := resp.Header.Get("Content-Length")
	if ct == "" || cl == "" {
		return f, errors.New("can't get details content-type: " + ct + ", content-length: " + cl)
	}
	f.ContentType = ct

	clParsed, err := strconv.ParseInt(cl, 10, 32)
	if err != nil {
		l.WithError(err).Errorf("can't parse content-length: " + cl)
		return f, err
	}

	f.ContentLength = clParsed
	return f, nil
}

func (f FileDetails) GetCache(url *url.URL, id string) (FileDetails, error) {
	cacheKey := cacheKeyPrefix + id
	_, err := cache.Client.GetKey(cacheKey, &f)
	if err == nil {
		return f, nil
	}

	f, err = f.Get(url)
	if err != nil {
		return f, err
	}

	go cache.Client.MarshalAndSetKey(cacheKey, f, cacheTTL)

	return f, nil
}

type FileDetailsDeps interface {
	Get(url.URL) (FileDetails, error)
	GetCache(url.URL, string) (FileDetails, error)
}

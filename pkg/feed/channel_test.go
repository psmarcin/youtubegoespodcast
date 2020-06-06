package feed

import (
	"github.com/eduncan911/podcast"
	"testing"
	"time"
)

func _TestFeed_addItem(t *testing.T) {
	f,_ := Create("123")
	ti := time.Now()
	f.Content = podcast.New("title", "http://onet", "description", &ti, &ti)
}

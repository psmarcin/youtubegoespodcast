package feed

import (
	"github.com/eduncan911/podcast"
	"strings"
	"testing"
	"time"
)

func TestFeed_serialize(t *testing.T) {
	f := New("123")
	ti := time.Now()
	f.Content = podcast.New("title", "http://onet", "description", &ti, &ti)
	serialized, err := f.Serialize()
	if err != nil {
		t.Errorf("serialize should not return error")
	}

	ok := strings.Contains(string(serialized), "<link>http://onet</link>")

	if !ok {
		t.Errorf("should contain link but it doesn't")
	}
}

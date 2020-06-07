package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"

	"github.com/psmarcin/youtubegoespodcast/pkg/config"
)

// Cache keeps firestor client and colletion
type Cache struct {
	firestore  *firestore.Client
	collection *firestore.CollectionRef
}

// Entity is a model for Firestore that handles cache with ttl
type Entity struct {
	Key   string        `firestore:"key"`
	Value string        `firestore:"value"`
	Ttl   time.Duration `firestore:"ttl"`
}

var Client Cache
var l = logrus.WithField("source", "cache")

// SetKey saves value for key in cache store with expiration time
func (c *Cache) SetKey(key, value string, exp time.Duration) error {
	// todo: pass external context
	ctx := context.Background()

	_, err := c.collection.Doc(key).Set(ctx, Entity{
		Key:   key,
		Value: value,
		Ttl:   exp,
	})
	if err != nil {
		l.WithError(err).WithFields(logrus.Fields{
			"key":   key,
			"value": value,
			"exp":   exp.String(),
		}).Errorf("set failed for %s", key)
		return err
	}

	return nil
}

// MarshalAndSetKey is the same as SetKey but before that it marshals value. Simple helper
func (c *Cache) MarshalAndSetKey(key string, value interface{}, exp time.Duration) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		l.WithError(err).WithField("value", value).WithField("key", key).Errorf("can't marshal")
		return err
	}
	err = c.SetKey(key, string(marshaled), exp)
	if err != nil {
		l.WithError(err).Error("can't set marshaled value")
	}
	return nil
}

// GetKey retrieve value by key from cache store
func (c *Cache) GetKey(key string, to interface{}) (string, error) {
	// todo: pass external context
	ctx := context.Background()
	raw, err := c.collection.Doc(key).Get(ctx)
	if err != nil {
		l.WithError(err).Warnf("can't get document %s", key)
		return "", err
	}

	var e Entity
	err = raw.DataTo(&e)
	if err != nil {
		l.WithError(err).Errorf("can't parse document %s", key)
		return "", err
	}

	timeDiff := time.Now().Sub(raw.UpdateTime)
	if timeDiff > e.Ttl {
		l.Debugf("key expires, updated at %s, expires in %s", raw.UpdateTime.Format(time.RFC3339), e.Ttl.String())
		return "", errors.New("Cache expires for " + key)
	}

	if to != nil {
		err = json.Unmarshal([]byte(e.Value), to)
	}

	if err != nil {
		return "", err
	}

	l.Debugf("got key %s", key)

	return e.Value, nil
}

// Connect establishes connection to Firebase and update singleton variable
func Connect() (Cache, error) {
	cache := Cache{}
	ctx := context.Background()
	store, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		logrus.WithError(err).Errorf("can't connect to Firebase")
		return cache, err
	}
	//defer store.Close() // Close client when done.

	cache.firestore = store
	cache.collection = store.Collection(config.Cfg.FirestoreCollection)
	Client = cache

	return cache, nil
}

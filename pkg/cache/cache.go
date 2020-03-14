package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sirupsen/logrus"

	"ygp/pkg/config"
)

type Cache struct {
	firestore  *firestore.Client
	collection *firestore.CollectionRef
}

type Entity struct {
	Key   string `firestore:"key"`
	Value string `firestore:"value"`
	Ttl time.Duration `firestore:"ttl"`
}

var Client Cache

func (c *Cache) SetKey(key, value string, exp time.Duration) error {
	// todo: pass external context
	ctx := context.Background()

	document, err := c.collection.Doc(key).Set(ctx, Entity{
		Key:   key,
		Value: value,
		Ttl: exp,
	})
	if err != nil {
		logrus.WithError(err).Fatalf("[CACHE] Set failed for %s", key)
		return err
	}

	logrus.WithField("document", document).Debugf("Set key %s", key)

	return nil
}

func (c *Cache) GetKey(key string, to interface{}) (string, error) {
	// todo: pass external context
	ctx := context.Background()
	raw, err := c.collection.Doc(key).Get(ctx)
	if err != nil {
		logrus.WithError(err).Errorf("[CACHE] Can't get document %s", key)
		return "", err
	}

	var e Entity
	err = raw.DataTo(&e)
	if err != nil{
		logrus.WithError(err).Errorf("[CACHE] Can't parse document %s", key)
		return "", err
	}

	timeDiff := time.Now().Sub(raw.UpdateTime)
	if timeDiff > e.Ttl {
		logrus.Debugf("[CACHE] Key expires, updated at %s, expires in %s", raw.UpdateTime.Format(time.RFC3339), e.Ttl.String())
		return "", errors.New("Cache expires for " + key)
	}

	if to != nil {
		err = json.Unmarshal([]byte(e.Value), to)
	}

	if err != nil {
		return "", err
	}

	logrus.Debugf("[CACHE] Got key %s", key)

	return e.Value, nil
}

func Connect() (Cache, error) {
	cache := Cache{}
	ctx := context.Background()
	store, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		logrus.WithError(err).Errorf("[CACHE] Can't connect to Firebase")
		return cache, err
	}
	// todo: investigate problem with disconnection
	//defer store.Close() // Close client when done.
	cache.firestore = store
	cache.collection = store.Collection(config.Cfg.FirestoreCollection)
	Client = cache
	return cache, nil
}

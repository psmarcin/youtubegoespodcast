package adapters

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"

	"cloud.google.com/go/firestore"
	"github.com/psmarcin/youtubegoespodcast/internal/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

// Cache keeps firestore client and collection
type Cache struct {
	firestore  *firestore.Client
	collection *firestore.CollectionRef
}

// CacheEntity is a model for Firestore that handles cache with ttl
type CacheEntity struct {
	Key   string        `firestore:"key"`
	Value string        `firestore:"value"`
	Ttl   time.Duration `firestore:"ttl"`
}

var l = logrus.WithField("source", "adapter")
var tracer = otel.GetTracerProvider().Tracer("yt.psmarcin.dev/adapters")

// NewCacheRepository establishes connection to Firebase and update singleton variable
func NewCacheRepository() (Cache, error) {
	cache := Cache{}
	ctx := context.Background()
	store, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		logrus.WithError(err).Errorf("can't connect to Firebase")
		return cache, errors.Wrap(err, "can't connect to Firebase")
	}

	cache.firestore = store
	cache.collection = store.Collection(config.Cfg.FirestoreCollection)

	return cache, nil
}

// SetKey saves value for key in cache store with expiration time
func (c *Cache) SetKey(ctx context.Context, key, value string, exp time.Duration) error {
	ctx, span := tracer.Start(ctx, "set-key")
	span.SetAttributes(label.String("key", key))
	span.SetAttributes(label.String("value", value))
	defer span.End()
	_, err := c.collection.Doc(key).Set(ctx, CacheEntity{
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
		span.RecordError(err)
		return errors.Wrapf(err, "can't set document: %s with value: %s, ctx: %s", key, value, ctx)
	}

	return nil
}

// GetKey retrieve value by key from cache store
func (c *Cache) GetKey(ctx context.Context, key string) (string, error) {
	ctx, span := tracer.Start(ctx, "get-key")
	span.SetAttributes(label.String("key", key))
	defer span.End()

	raw, err := c.collection.Doc(key).Get(ctx)
	if err != nil {
		l.WithError(err).Warnf("can't get document %s", key)
		span.RecordError(err)
		return "", errors.Wrapf(err, "can't get document from database key: %s with context %s", key, ctx)
	}

	var e CacheEntity
	err = raw.DataTo(&e)
	if err != nil {
		l.WithError(err).Errorf("can't parse document %s", key)
		span.RecordError(err)
		return "", errors.Wrap(err, "can't parse raw data")
	}

	timeDiff := time.Since(raw.UpdateTime)
	if timeDiff > e.Ttl {
		l.Debugf("key expires, updated at %s, expires in %s", raw.UpdateTime.Format(time.RFC3339), e.Ttl.String())
		span.RecordError(err)
		return "", errors.Newf("Cache expires for %s", key)
	}

	l.Debugf("got key %s", key)

	return e.Value, nil
}

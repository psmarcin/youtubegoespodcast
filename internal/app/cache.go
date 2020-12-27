package app

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"time"
)

const (
	CacheTTL = time.Hour * 24
)

var l = logrus.WithField("source", "app")
var tracer = otel.GetTracerProvider().Tracer("yt.psmarcin.dev/app")

type cacheRepository interface {
	SetKey(context.Context, string, string, time.Duration) error
	GetKey(context.Context, string) (string, error)
}

type CacheService struct {
	cache cacheRepository
}

func NewCacheService(cache cacheRepository) CacheService {
	if cache == nil {
		panic("missing cacheRepository")
	}

	return CacheService{
		cache: cache,
	}
}

// Set is the same as SetKey but before that it marshals value. Simple helper
func (c *CacheService) Set(ctx context.Context, key string, value interface{}) error {
	ctx, span := tracer.Start(ctx, "set")
	span.SetAttributes(label.String("key", key))
	span.SetAttributes(label.Any("value", value))
	defer span.End()

	marshaled, err := json.Marshal(value)
	if err != nil {
		l.WithError(err).WithField("value", value).WithField("key", key).Errorf("can't marshal")
		span.RecordError(err)
		return err
	}
	err = c.cache.SetKey(ctx, key, string(marshaled), CacheTTL)
	if err != nil {
		l.WithError(err).Error("can't set marshaled value")
		span.RecordError(err)
	}
	return nil
}

// Get looks for object via key in cache
func (c *CacheService) Get(ctx context.Context, key string, to interface{}) error {
	tCtx, span := tracer.Start(ctx, "get")
	span.SetAttributes(label.Any("key", key))
	defer span.End()

	entity, err := c.cache.GetKey(tCtx, key)
	if err != nil {
		span.RecordError(err)
		return errors.Wrap(err, "can't get cache for CacheService")
	}

	if to != nil {
		err = json.Unmarshal([]byte(entity), to)
	}

	if err != nil {
		span.RecordError(err)
		return errors.Wrapf(err, "can't unmarshal value for key %s", key)
	}

	return nil
}

// MarshalAndSetKey is the same as SetKey but before that it marshals value. Simple helper
func (c *CacheService) MarshalAndSet(ctx context.Context, key string, value interface{}) error {
	tCtx, span := tracer.Start(ctx, "marshal-and-set")
	span.SetAttributes(label.String("key", key))
	span.SetAttributes(label.Any("value", value))
	defer span.End()

	marshaled, err := json.Marshal(value)
	if err != nil {
		l.WithError(err).WithField("value", value).WithField("key", key).Errorf("can't marshal")
		span.RecordError(err)
		return err
	}
	err = c.Set(tCtx, key, string(marshaled))
	if err != nil {
		l.WithError(err).Error("can't set marshaled value")
		span.RecordError(err)
	}
	return nil
}

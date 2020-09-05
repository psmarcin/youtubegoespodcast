package app

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	CacheTTL = time.Hour * 24
)

var l = logrus.WithField("source", "app")

type cacheRepository interface {
	SetKey(string, string, time.Duration) error
	GetKey(string) (string, error)
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
func (c *CacheService) Set(key string, value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		l.WithError(err).WithField("value", value).WithField("key", key).Errorf("can't marshal")
		return err
	}
	err = c.cache.SetKey(key, string(marshaled), CacheTTL)
	if err != nil {
		l.WithError(err).Error("can't set marshaled value")
	}
	return nil
}

// Get looks for object via key in cache
func (c *CacheService) Get(key string, to interface{}) error {
	entity, err := c.cache.GetKey(key)
	if err != nil {
		return errors.Wrap(err, "can't get cache for CacheService")
	}

	if to != nil {
		err = json.Unmarshal([]byte(entity), to)
	}

	if err != nil {
		return errors.Wrapf(err, "can't unmarshal value for key %s", key)
	}

	return nil
}

// MarshalAndSetKey is the same as SetKey but before that it marshals value. Simple helper
func (c *CacheService) MarshalAndSet(key string, value interface{}) error {
	marshaled, err := json.Marshal(value)
	if err != nil {
		l.WithError(err).WithField("value", value).WithField("key", key).Errorf("can't marshal")
		return err
	}
	err = c.Set(key, string(marshaled))
	if err != nil {
		l.WithError(err).Error("can't set marshaled value")
	}
	return nil
}

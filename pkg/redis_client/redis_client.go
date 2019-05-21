package redis_client

import (
	"encoding/json"
	"errors"
	"time"
	"ytg/pkg/config"

	"github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
)

var (
	errNotConnected = errors.New("Redis client not connected")
)

type redisClient struct {
	client    *redis.Client
	connected bool
}

// Client hold all informations and instance to do queries
var Client redisClient

func (r *redisClient) SetKey(key, value string, exp time.Duration) error {
	if !r.connected {
		return errNotConnected
	}
	logrus.Printf("[REDIS] About to SET %s", key)
	_, err := r.client.Set(key, value, exp).Result()
	return err
}

func (r *redisClient) GetKey(key string, to interface{}) (string, error) {
	if !r.connected {
		return "", errNotConnected
	}
	logrus.Printf("[REDIS] About to GET %s", key)
	val, err := r.client.Get(key).Result()
	if err != nil {
		return "", err
	}
	if to != nil {
		err = json.Unmarshal([]byte(val), to)
		if err != nil {
			return "", err
		}
	}
	return val, err
}

// Connect creates connection to redis base on config package
func Connect() redisClient {
	Client = redisClient{
		connected: false,
	}

	Client.client = redis.NewClient(&redis.Options{
		Addr: config.Cfg.RedisHost + ":" + config.Cfg.RedisPort,
	})
	Client.connected = true
	return Client
}

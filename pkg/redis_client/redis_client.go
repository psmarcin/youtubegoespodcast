package redis_client

import (
	"encoding/json"
	"errors"
	"time"
	"ygp/pkg/config"

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

// Client hold all information and instance to do queries
var Client redisClient

func (r *redisClient) SetKey(key, value string, exp time.Duration) error {
	if !r.connected {
		return errNotConnected
	}
	logrus.Printf("[CACHE] SET %s", key)
	_, err := r.client.Set(key, value, exp).Result()
	return err
}

func (r *redisClient) GetKey(key string, to interface{}) (string, error) {
	if !r.connected {
		return "", errNotConnected
	}
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
	opt, err := redis.ParseURL(config.Cfg.RedisURI)
	if err != nil {
		panic(err)
	}

	Client = redisClient{
		connected: false,
	}

	Client.client = redis.NewClient(opt)
	Client.connected = true
	logrus.Printf("[CACHE] Connected to %s", config.Cfg.RedisURI)
	return Client
}

// Teardown closes connection and cleanup after cache connection
func Teardown() {
	Client.client.Close()
	Client.connected = false
}

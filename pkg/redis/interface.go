package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Interface interface {
	Keys(pattern string) ([]string, error)

	Get(key string) (string, error)

	Set(key string, value string, duration time.Duration) error

	Del(keys ...string) error

	LPush(key string, values ...interface{}) error

	RPop(key string) (string, error)

	BRPop(keys ...string) ([]string, error)

	Exists(keys ...string) (bool, error)

	Expire(key string, duration time.Duration) error

	Subscribe(channels ...string) *redis.PubSub

	Publish(channel string, message interface{}) error
}

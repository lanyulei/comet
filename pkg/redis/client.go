package redis

import (
	"context"
	"fmt"
	"github.com/lanyulei/comet/pkg/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var (
	rc          Interface
	stopChRedis chan struct{}
)

type Client struct {
	client *redis.Client
}

func Setup() {
	var err error
	stopChRedis = make(chan struct{})
	rc, err = newRedisClient(
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.password"),
		viper.GetInt("redis.db"),
		stopChRedis,
	)

	if err != nil {
		logger.Fatal("redis initialization failed")
	}

	return
}

func newRedisClient(host string, port int, password string, db int, stopCh <-chan struct{}) (Interface, error) {
	var r Client

	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,

		PoolSize:     20,
		MinIdleConns: 10,
	}

	if stopCh == nil {
		logger.Fatalf("no stop channel passed, redis connections will leak.")
	}

	r.client = redis.NewClient(redisOptions)

	if err := r.client.Ping(context.TODO()).Err(); err != nil {
		_ = r.client.Close()
		return nil, err
	}

	// close redis in case of connection leak
	if stopCh != nil {
		go func() {
			<-stopCh
			if err := r.client.Close(); err != nil {
				logger.Error(err)
			} else {
				logger.Info("redis connection closed")
			}
		}()
	}

	return &r, nil
}

func (r *Client) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

func (r *Client) Keys(pattern string) ([]string, error) {
	return r.client.Keys(context.TODO(), pattern).Result()
}

func (r *Client) Set(key string, value string, duration time.Duration) error {
	return r.client.Set(context.TODO(), key, value, duration).Err()
}

func (r *Client) Del(keys ...string) error {
	return r.client.Del(context.TODO(), keys...).Err()
}

func (r *Client) LPush(key string, values ...interface{}) error {
	return r.client.LPush(context.TODO(), key, values...).Err()
}

func (r *Client) RPop(key string) (string, error) {
	return r.client.RPop(context.TODO(), key).Result()
}

func (r *Client) BRPop(keys ...string) ([]string, error) {
	return r.client.BRPop(context.TODO(), time.Duration(0), keys...).Result()
}

func (r *Client) Exists(keys ...string) (bool, error) {
	existedKeys, err := r.client.Exists(context.TODO(), keys...).Result()
	if err != nil {
		return false, err
	}

	return len(keys) == int(existedKeys), nil
}

func (r *Client) Expire(key string, duration time.Duration) error {
	return r.client.Expire(context.TODO(), key, duration).Err()
}

func (r *Client) Subscribe(channels ...string) *redis.PubSub {
	return r.client.Subscribe(context.TODO(), channels...)
}

func (r *Client) Publish(channel string, message interface{}) error {
	return r.client.Publish(context.TODO(), channel, message).Err()
}

func Rc() Interface {
	return rc
}

func StopChRedis() chan struct{} {
	return stopChRedis
}

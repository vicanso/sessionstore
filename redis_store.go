package sessionstore

import (
	"errors"

	"github.com/go-redis/redis"
)

type (
	// RedisStore redis store for session
	RedisStore struct {
		client *redis.Client
		Store
	}
)

// Get get the session from redis
func (rs *RedisStore) Get(key string) ([]byte, error) {
	buf, err := rs.client.Get(key).Bytes()
	if err == redis.Nil {
		return buf, nil
	}
	return buf, err
}

// Set set the session to redis
func (rs *RedisStore) Set(key string, data []byte) error {
	return rs.client.Set(key, data, rs.GetTTL()).Err()
}

// Destroy remove the session from redis
func (rs *RedisStore) Destroy(key string) error {
	return rs.client.Del(key).Err()
}

// NewRedisClient create a new redis client
func NewRedisClient(opts *redis.Options) *redis.Client {
	return redis.NewClient(opts)
}

// NewRedisStore create new redis store instance
func NewRedisStore(client *redis.Client) *RedisStore {
	if client == nil {
		panic(errors.New("client can not be nil"))
	}
	rs := &RedisStore{
		client: client,
	}
	return rs
}

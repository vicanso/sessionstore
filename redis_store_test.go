package sessionstore

import (
	"bytes"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestRedisStore(t *testing.T) {
	key := generateID()
	data := []byte("tree.xie")
	ttl := 300 * time.Second
	var rs *RedisStore

	t.Run("new redis store", func(t *testing.T) {
		client := NewRedisClient(&redis.Options{
			Addr: "localhost:6379",
		})
		rs = NewRedisStore(client)
		rs.SetOptions(&Options{
			TTL: ttl,
		})
	})
	t.Run("get not exists data", func(t *testing.T) {
		buf, err := rs.Get(key)
		if err != nil || len(buf) != 0 {
			t.Fatalf("should return empty bytes")
		}
	})

	t.Run("set data", func(t *testing.T) {
		err := rs.Set(key, data)
		if err != nil {
			t.Fatalf("set data fail, %v", err)
		}
		buf, err := rs.Get(key)
		if err != nil {
			t.Fatalf("get data fail after set, %v", err)
		}
		if !bytes.Equal(data, buf) {
			t.Fatalf("the data is not the same after set")
		}
	})

	t.Run("destroy", func(t *testing.T) {
		err := rs.Destroy(key)
		if err != nil {
			t.Fatalf("destory data fail, %v", err)
		}
		buf, err := rs.Get(key)
		if err != nil || len(buf) != 0 {
			t.Fatalf("shoud return empty bytes after destroy")
		}
	})
}

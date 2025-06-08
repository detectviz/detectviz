package redis_test

import (
	"testing"
	"time"

	"github.com/detectviz/detectviz/internal/adapters/cachestore/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newTestClient() *redis.RedisCacheStore {
	client := redis.NewRedisCacheStore(redisclient())
	_ = client.Delete("test:key1")
	_ = client.Delete("test:key2")
	_ = client.Delete("test:temp")
	return client
}

func redisclient() *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}

func TestRedisCacheStore_CRUD(t *testing.T) {
	store := newTestClient()

	err := store.Set("test:key1", "value1", 10)
	assert.NoError(t, err)

	has, err := store.Has("test:key1")
	assert.NoError(t, err)
	assert.True(t, has)

	val, err := store.Get("test:key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	err = store.Delete("test:key1")
	assert.NoError(t, err)

	has, err = store.Has("test:key1")
	assert.NoError(t, err)
	assert.False(t, has)

	val, err = store.Get("test:key1")
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestRedisCacheStore_TTL(t *testing.T) {
	store := newTestClient()

	err := store.Set("test:temp", "soon-expire", 1)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	has, err := store.Has("test:temp")
	assert.NoError(t, err)
	assert.False(t, has)

	val, err := store.Get("test:temp")
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestRedisCacheStore_Keys(t *testing.T) {
	store := newTestClient()

	err := store.Set("test:key1", "1", 10)
	assert.NoError(t, err)
	err = store.Set("test:key2", "2", 10)
	assert.NoError(t, err)

	keys, err := store.Keys("test:")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"test:key1", "test:key2"}, keys)
}

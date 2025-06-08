package cachestore_test

import (
	"testing"

	"github.com/detectviz/detectviz/internal/adapters/cachestore"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetMemoryBackend(t *testing.T) {
	store, err := cachestore.Get("memory")
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestGetUnknownBackend(t *testing.T) {
	store, err := cachestore.Get("unknown")
	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestWithRedisClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})
	cachestore.WithRedisClient(client)

	store, err := cachestore.Get("redis")
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestListBackends(t *testing.T) {
	names := cachestore.List()
	assert.Contains(t, names, "memory")
	assert.Contains(t, names, "redis")
}

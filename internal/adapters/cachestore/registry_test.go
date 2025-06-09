package cachestoreadapter_test

import (
	"testing"

	cachestoreadapter "github.com/detectviz/detectviz/internal/adapters/cachestore"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetMemoryBackend(t *testing.T) {
	store, err := cachestoreadapter.Get("memory")
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestGetUnknownBackend(t *testing.T) {
	store, err := cachestoreadapter.Get("unknown")
	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestWithRedisClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})
	cachestoreadapter.WithRedisClient(client)

	store, err := cachestoreadapter.Get("redis")
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestListBackends(t *testing.T) {
	names := cachestoreadapter.List()
	assert.Contains(t, names, "memory")
	assert.Contains(t, names, "redis")
}

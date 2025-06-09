package memoryadapter_test

import (
	"testing"

	memoryadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCacheStore_BasicOperations(t *testing.T) {
	store := memoryadapter.NewMemoryCacheStore()

	err := store.Set("foo", "123", 0)
	assert.NoError(t, err)
	err = store.Set("bar", "456", 0)
	assert.NoError(t, err)

	has, err := store.Has("foo")
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = store.Has("bar")
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = store.Has("baz")
	assert.NoError(t, err)
	assert.False(t, has)

	val, err := store.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "123", val)

	val, err = store.Get("bar")
	assert.NoError(t, err)
	assert.Equal(t, "456", val)

	val, err = store.Get("baz")
	assert.NoError(t, err)
	assert.Equal(t, "", val)

	err = store.Delete("foo")
	assert.NoError(t, err)

	has, err = store.Has("foo")
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestMemoryCacheStore_Keys(t *testing.T) {
	store := memoryadapter.NewMemoryCacheStore()

	err := store.Set("user:1", "A", 0)
	assert.NoError(t, err)
	err = store.Set("user:2", "B", 0)
	assert.NoError(t, err)
	err = store.Set("device:1", "X", 0)
	assert.NoError(t, err)

	keys, err := store.Keys("user:")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"user:1", "user:2"}, keys)
}

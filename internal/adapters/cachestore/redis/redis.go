package redisadapter

import (
	"context"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// RedisCacheStore 是 Redis 實作的 CacheStore。
// zh: 使用 Redis 作為儲存後端，支援快取存取與過期控制。
type RedisCacheStore struct {
	client *goredis.Client
	ctx    context.Context
}

// NewRedisCacheStore 建立 Redis 快取實例。
// zh: 須傳入 go-redis v9 的 Redis 客戶端。
func NewRedisCacheStore(client *goredis.Client) *RedisCacheStore {
	return &RedisCacheStore{
		client: client,
		ctx:    context.Background(),
	}
}

// Get 取得指定 key 的值。
// zh: 若 key 不存在則回傳空字串與 nil 錯誤。
func (r *RedisCacheStore) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == goredis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set 寫入 key 的值與過期時間（秒）。
// zh: ttl 以秒為單位，若 ttl <= 0 則為永久存留。
func (r *RedisCacheStore) Set(key, value string, ttl int) error {
	return r.client.Set(r.ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

// Has 檢查 key 是否存在。
// zh: 若存在回傳 true，否則回傳 false，若查詢失敗則回傳 error。
func (r *RedisCacheStore) Has(key string) (bool, error) {
	val, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// Delete 移除指定 key。
// zh: 若 key 存在則刪除，否則無動作，回傳刪除錯誤（若有）。
func (r *RedisCacheStore) Delete(key string) error {
	_, err := r.client.Del(r.ctx, key).Result()
	return err
}

// Keys 回傳符合 prefix 的所有 key。
// zh: 僅回傳以 prefix 開頭的 key 清單。
func (r *RedisCacheStore) Keys(prefix string) ([]string, error) {
	keys, err := r.client.Keys(r.ctx, prefix+"*").Result()
	if err != nil {
		return nil, err
	}
	var filtered []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			filtered = append(filtered, k)
		}
	}
	return filtered, nil
}

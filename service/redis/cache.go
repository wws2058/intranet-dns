package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tswcbyy1107/intranet-dns/database"
)

// get cache
func LoadCache(key string, dst interface{}) error {
	return GetToDst(key, dst)
}

// cache, k,v,ttl
func Cache(key string, value interface{}, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	bytes, _ := json.Marshal(value)
	_, err := database.Rdb.Set(ctx, key, string(bytes), ttl).Result()
	return err
}

// del cache
func DropCache(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	database.Rdb.Del(ctx, key)
}

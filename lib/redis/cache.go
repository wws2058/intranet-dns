package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tswcbyy1107/dns-service/database"
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
	result, err := database.Rdb.Set(ctx, key, string(bytes), ttl).Result()
	fmt.Println("cache:", result, err)
	return err
}

// del cache
func DropCache(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	database.Rdb.Del(ctx, key)
}

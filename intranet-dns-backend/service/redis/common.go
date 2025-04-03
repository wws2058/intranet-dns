package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tswcbyy1107/intranet-dns/database"
)

// k v string, json.Unmarshal to dst, dst is pointer
func GetToDst(key string, dst interface{}) (err error) {
	value, err := DoCmd("GET", key)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value.(string)), dst)
	return
}

// redis exec cmd, 500ms timeout, key does not exist return "", args eg: set xx yy, slice
func DoCmd(args ...interface{}) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	value, err := database.Rdb.Do(ctx, args...).Result()
	// key does not exists
	if err != nil && err == redis.Nil {
		return "", nil
	}
	return value, err
}

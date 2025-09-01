package redis

import (
	"context"
	"time"

	"github.com/wws2058/intranet-dns/database"
	"github.com/wws2058/intranet-dns/utils"
)

// require defer UnLock
func Lock(key string, ttl time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	value := time.Now().Format(utils.DefaultTimeFormat)
	ok := database.Rdb.SetNX(ctx, key, value, ttl).Val()
	return ok
}

func UnlockNow(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	database.Rdb.Del(ctx, key)
}

// try lock every 50ms in retry period
func TryLock(key string, ttl, retry time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), retry)
	defer cancel()

	value := time.Now().Format(utils.DefaultTimeFormat)
	ok := database.Rdb.SetNX(ctx, key, value, ttl).Val()
	if ok {
		return true
	}
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			ok := database.Rdb.SetNX(ctx, key, value, ttl).Val()
			if ok {
				return true
			}
		}
	}
}

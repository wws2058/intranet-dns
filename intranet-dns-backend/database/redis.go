package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/wws2058/intranet-dns/config"
)

var Rdb *redis.Client

func initRdb() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%v", config.GlobalConfig.Redis.Host, config.GlobalConfig.Redis.Port),
		Password:        config.GlobalConfig.Redis.Passwd,
		DB:              config.GlobalConfig.Redis.Database,
		PoolSize:        100,
		MaxIdleConns:    10,
		ConnMaxIdleTime: time.Hour,
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.WithField("redis", "init").Error(err)
		os.Exit(1)
	}
	logrus.WithField("redis", "init").Info("succeed")
}

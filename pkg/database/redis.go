package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/ilooky/go-layout/pkg/config"
	"github.com/ilooky/logger"
	"strconv"
	"time"
)

func InitRedis(conf config.Redis) {
	var ctx = context.Background()
	dbIndex, _ := strconv.Atoi(conf.Database)
	ring := redis.NewClient(&redis.Options{
		Addr:         conf.Host + ":" + conf.Port,
		Password:     conf.Password,
		DB:           dbIndex,
		DialTimeout:  30 * time.Second,
		MinIdleConns: 3,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	if _, err := ring.Ping(ctx).Result(); err != nil {
		logger.Debug(err)
		return
	}
}

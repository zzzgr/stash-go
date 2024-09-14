package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

var Rdb *redis.Client
var ctx = context.Background()

// Init redis
func Init(addr, username, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	ping := Rdb.Ping(context.Background())
	if err := ping.Err(); err != nil {
		log.Fatalf("redis连接失败, 请检查 (%s)", err.Error())
	}
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return Rdb.Set(ctx, key, value, expiration).Err()
}

func LPush(key string, value interface{}, expiration time.Duration) error {
	cmd := Rdb.LPush(ctx, key, value)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	c := Rdb.Expire(ctx, key, expiration)
	return c.Err()
}

func List(key string) []string {
	l := Rdb.LLen(ctx, key)

	cmd := Rdb.LRange(ctx, key, 0, l.Val())
	return cmd.Val()
}

func Del(keys ...string) error {
	return Rdb.Del(ctx, keys...).Err()
}

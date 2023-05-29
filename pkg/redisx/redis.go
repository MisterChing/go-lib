package redisx

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr, password string, db ...int) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	}
	if len(db) > 0 {
		opt.DB = db[0]
	}
	redisIns := redis.NewClient(opt)
	ret, err := redisIns.Ping(context.Background()).Result()
	if err != nil || ret != "PONG" {
		return nil, fmt.Errorf("connect redis failed. error:%+v", err)
	}
	return redisIns, nil
}

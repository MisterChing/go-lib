package cachex

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

func RedisMGet(ctx context.Context, redisIns *redis.Client, keyPrefix string, entityIds []string) (hitData map[string]interface{}, missedIds []string, err error) {
	var (
		keySlice []string
		retSlice []interface{}
	)
	hitData = make(map[string]interface{})
	keyPrefix = strings.Trim(keyPrefix, ":")
	for _, v := range entityIds {
		tmpKey := fmt.Sprintf("%s:%s", keyPrefix, v)
		keySlice = append(keySlice, tmpKey)
	}
	retSlice, err = redisIns.MGet(ctx, keySlice...).Result()
	if err != nil {
		return
	}
	for k, id := range entityIds {
		if retSlice[k] == "" {
			missedIds = append(missedIds, id)
		} else {
			hitData[id] = retSlice[k]
		}
	}
	return
}

func RedisMSet(ctx context.Context, redisIns *redis.Client, keyPrefix string, data map[string]interface{}, expire time.Duration) error {
	keyPrefix = strings.Trim(keyPrefix, ":")
	for k, v := range data {
		tmpKey := fmt.Sprintf("%s:%s", keyPrefix, k)
		if err := redisIns.SetEX(ctx, tmpKey, v, expire).Err(); err != nil {
			return err
		}
	}
	return nil
}

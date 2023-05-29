package redisx

import (
	"context"
	"log"
	"testing"

	"github.com/MisterChing/go-lib/utils/debugutil"
)

func TestRedis(t *testing.T) {

	redis, err := NewRedis("xxxx:6379", "xxxx", 0)
	if err != nil {
		log.Fatal(err)
	}
	iter := redis.Scan(context.Background(), 0, "ching*", 1000).Iterator()
	for iter.Next(context.Background()) {
		if iter.Err() != nil {
			log.Fatal(iter.Err())
			return
		}
		debugutil.DebugPrint(iter.Val(), 0)
	}
}

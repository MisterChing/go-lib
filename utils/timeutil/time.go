package timeutil

import (
	"log"
	"time"
)

func TimeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		log.Println("time cost:", tc)
	}
}

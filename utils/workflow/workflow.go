package workflow

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"sync"
	"time"
)

// GoWithNumber wg不能复制 可以用vet 检测 go vet -copylocks ./xxx.go
func GoWithNumber(wg *sync.WaitGroup, goroutineNum int, fn func()) {
	wg.Add(goroutineNum)
	for i := 0; i < goroutineNum; i++ {
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
}

func GoWithRecover(c *gin.Context, fn func()) {
	cp := c.Copy()
	go func() {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				if c != nil {
					fmt.Println(string(debug.Stack()))
				} else {
					fmt.Println(string(debug.Stack()))
				}
			}
		}(cp)
		fn()
	}()
}

//func GoWithRecover(fn func()) {
//	go func() {
//		defer func() {
//			if err := recover(); err != nil {
//				fmt.Println(string(debug.Stack()))
//			}
//		}()
//		fn()
//	}()
//}

type FuncWithArgs func(args ...interface{})

func GoWithRecoverWithArgs(c *gin.Context, fn FuncWithArgs, args ...interface{}) {
	cp := c.Copy()
	go func() {
		defer func(c *gin.Context) {
			if err := recover(); err != nil {
				if c != nil {
					fmt.Println(string(debug.Stack()))
				} else {
					fmt.Println(string(debug.Stack()))
				}
			}
		}(cp)
		fn(args...)
	}()
}

func GoGroupWait(c *gin.Context, fnArr ...func() string) {
	if len(fnArr) == 0 {
		return
	}
	var wg sync.WaitGroup
	goNum := len(fnArr)
	wg.Add(goNum)
	for _, fn := range fnArr {
		cp := c.Copy()
		go func(doFn func() string) {
			defer func(c *gin.Context) {
				if err := recover(); err != nil {
					if c != nil {
						fmt.Println(string(debug.Stack()))
					} else {
						fmt.Println(string(debug.Stack()))
					}
				}
			}(cp)
			defer wg.Done()
			doFn()
		}(fn)
	}
	wg.Wait()
}

func WithRetry(fn func() error, retryCount int) {
	for ; retryCount > 0; retryCount-- {
		if err := fn(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}

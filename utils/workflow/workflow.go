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

func GinGoWithRecover(c *gin.Context, fn func(c *gin.Context)) {
	cp := c.Copy()
	go func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if c != nil {
					fmt.Println(string(debug.Stack()))
				} else {
					fmt.Println(string(debug.Stack()))
				}
			}
		}()
		fn(c)
	}(cp)
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

type FuncWithArgs func(c *gin.Context, args ...interface{})

func GinGoWithRecoverWithArgs(c *gin.Context, fn FuncWithArgs, args ...interface{}) {
	cp := c.Copy()
	go func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if c != nil {
					fmt.Println(string(debug.Stack()))
				} else {
					fmt.Println(string(debug.Stack()))
				}
			}
		}()
		fn(cp, args...)
	}(cp)
}

func GinGoGroupWait(c *gin.Context, fnArr ...func(c *gin.Context) string) {
	if len(fnArr) == 0 {
		return
	}
	var wg sync.WaitGroup
	goNum := len(fnArr)
	wg.Add(goNum)
	for _, fn := range fnArr {
		cp := c.Copy()
		doFn := fn
		go func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					if c != nil {
						fmt.Println(string(debug.Stack()))
					} else {
						fmt.Println(string(debug.Stack()))
					}
				}
			}()
			defer wg.Done()
			doFn(c)
		}(cp)
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

package interval_manager

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"sync/atomic"
	"syscall"
	"time"
)

type IntervalManagerBizFn interface {
	BusinessFn(ctx context.Context, buf <-chan interface{}) error
	OnShutDownFn(ctx context.Context, buf <-chan interface{}) error
}

type IntervalManager struct {
	name         string
	tickDu       time.Duration
	buf          chan interface{}
	bufSize      int
	quitCh       chan os.Signal
	workerQuitCh chan struct{}
	bizFn        IntervalManagerBizFn
	successQuit  chan struct{}
	runningNum   int32
}

func NewIntervalManager(name string, qSize int, tickDu time.Duration) *IntervalManager {
	obj := &IntervalManager{
		name:         name,
		quitCh:       make(chan os.Signal, 1),
		buf:          make(chan interface{}, qSize),
		bufSize:      qSize,
		workerQuitCh: make(chan struct{}),
		successQuit:  make(chan struct{}, 1),
		tickDu:       tickDu,
	}
	return obj
}

func (manager *IntervalManager) SetBizFn(fn IntervalManagerBizFn) {
	manager.bizFn = fn
}

func (manager *IntervalManager) GetName() string {
	return manager.name
}

func (manager *IntervalManager) intervalWrap() {
	ticker := time.NewTicker(manager.tickDu)
	limitQ := make(chan struct{}, 1) //同一时刻只能有一个实例执行
	for {
		select {
		case <-ticker.C:
			limitQ <- struct{}{}
			atomic.AddInt32(&manager.runningNum, 1)
			ctx := genGinContext()
			go func() {
				defer func() {
					if err := recover(); err != nil {
                        //todo log
					}
				}()
				defer func() {
					atomic.AddInt32(&manager.runningNum, -1)
					<-limitQ
				}()
				//启动定时消费
				_ = manager.bizFn.BusinessFn(ctx, manager.buf)
			}()
		case <-manager.workerQuitCh:
			goto END
		}
	}
END:
	ticker.Stop()
	close(manager.buf)
	log.Printf("IntervalManager [%s] ticker stopped.", manager.name)
}

func (manager *IntervalManager) Produce(ctx context.Context, data interface{}) {
	if len(manager.buf) < manager.bufSize {
		manager.buf <- data
	} else {
		//立即消费
		_ = manager.bizFn.BusinessFn(ctx, manager.buf)
		manager.buf <- data
	}
}

func (manager *IntervalManager) Start() {
	log.Printf("IntervalManager [%s] start", manager.name)
	go manager.intervalWrap()
	//go func() {
	//	aTicker := time.NewTicker(time.Second)
	//	for {
	//		select {
	//		case <-aTicker.C:
	//			log.Println("name:", manager.name, "buf:", len(manager.buf), "running:", atomic.LoadInt32(&manager.runningNum))
	//		}
	//	}
	//}()

	// start monitor & block
	manager.monitor()
}

func (manager *IntervalManager) WaitStop() {
	<-manager.successQuit
}

func (manager *IntervalManager) monitor() {
	signal.Notify(manager.quitCh, syscall.SIGINT, syscall.SIGTERM)
	<-manager.quitCh
	log.Printf("received quit signal. IntervalManager [%s] start exiting...", manager.name)

	//通知worker退出
	close(manager.workerQuitCh)
	log.Printf("IntervalManager [%s] start waitting remained job...", manager.name)
	// 检测buf数据 & 启动收尾消费（限制最多5个协程）
	limitQ := make(chan struct{}, 5)
	for len(manager.buf) > 0 {
		limitQ <- struct{}{}
		atomic.AddInt32(&manager.runningNum, 1)
		go func() {
			ctx := genGinContext()
			defer func() {
				if err := recover(); err != nil {
                    //todo log
				}
			}()
			defer func() {
				atomic.AddInt32(&manager.runningNum, -1)
				<-limitQ
			}()
			_ = manager.bizFn.OnShutDownFn(ctx, manager.buf)
		}()
	}
	//等待全部running协程退出
	for atomic.LoadInt32(&manager.runningNum) > 0 {
		time.Sleep(500 * time.Millisecond)
	}
	log.Printf("IntervalManager [%s] all job stopped.", manager.name)
	manager.successQuit <- struct{}{}
	log.Printf("IntervalManager [%s] stopped.", manager.name)
}

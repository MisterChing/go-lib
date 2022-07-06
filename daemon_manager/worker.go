package daemon_manager

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type Worker interface {
	ProcessMessages(messages []interface{}) error
}

type WorkerManager struct {
	name         string
	bufSize      int
	workerMax    int
	stopDispatch chan struct{}
	graceQuit    chan struct{}
	quitCh       chan os.Signal
	runningNum   int32
	queue        *Queue
}

func NewWorkerManager(name string, bufSize int, workerMax int) *WorkerManager {
	obj := &WorkerManager{
		name:         name,
		bufSize:      bufSize,
		workerMax:    workerMax,
		stopDispatch: make(chan struct{}),
		graceQuit:    make(chan struct{}),
		quitCh:       make(chan os.Signal),
		queue:        NewQueue(bufSize),
	}
	return obj
}

func (manager *WorkerManager) Start() {
	log.Printf("DaemonManager [%s] start", manager.name)
	go manager.dispatch()

	// monitor & block
	manager.monitor()
}

func (manager *WorkerManager) WaitStop() {
	<-manager.graceQuit
}

func (manager *WorkerManager) Produce(data interface{}) {
	manager.queue.EnQueue(data)
}

func (manager *WorkerManager) dispatch() {

	for {
		select {
		case msg, ok := <-manager.queue.Buf():
			log.Println(ok, msg)
		case <-manager.stopDispatch:
			goto END
		}

	}
END:
	manager.queue.Close()
}

func (manager *WorkerManager) monitor() {
	manager.installSignals()
	<-manager.quitCh
	log.Printf("received quit signal. DaemonManager [%s] start exiting...", manager.name)

	// stop dispatch
	close(manager.stopDispatch)

	log.Printf("DaemonManager [%s] start waitting remained job...", manager.name)
	// check queue size
	limitQ := make(chan struct{}, 10)
	for manager.queue.Size() > 0 {
		limitQ <- struct{}{}
		atomic.AddInt32(&manager.runningNum, 1)
		msg := manager.queue.DeQueue()
		go func() {
			defer func() {
				atomic.AddInt32(&manager.runningNum, -1)
				<-limitQ
			}()
			log.Println(msg)
		}()
	}
	// waiting running worker
	for atomic.LoadInt32(&manager.runningNum) > 0 {
		time.Sleep(500 * time.Millisecond)
	}
	log.Printf("DaemonManager [%s] all job stopped.", manager.name)
	manager.graceQuit <- struct{}{}
	log.Printf("DaemonManager [%s] stopped.", manager.name)
}

func (manager *WorkerManager) installSignals() {
	signals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
	}
	signal.Notify(manager.quitCh, signals...)
}

package debugutil

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func DebugPrint(obj interface{}, isExit int) {
	if isExit == 1 {
		fmt.Printf("[debug2]---%p---%T---%+v\n", obj, obj, obj)
		os.Exit(0)
	} else {
		fmt.Printf("[debug2]---%p---%T---%+v\n", obj, obj, obj)
	}
}

func BlockMain(fnArr ...func()) {
	var (
		done = make(chan struct{})
		wg   sync.WaitGroup
	)
	go func() {
		quitCh := make(chan os.Signal)
		signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
		<-quitCh
		if len(fnArr) > 0 {
			for _, fn := range fnArr {
				wg.Add(1)
				f := fn
				go func() {
					defer wg.Done()
					f()
				}()
			}
			wg.Wait()
		}
		done <- struct{}{}
	}()
	<-done
}

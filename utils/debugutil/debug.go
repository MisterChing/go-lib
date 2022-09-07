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
		fmt.Printf("\033[1;40;32m[debug]\033[0m---%p---%T---%+v\n", obj, obj, obj)
		os.Exit(0)
	} else {
		fmt.Printf("\033[1;40;32m[debug]\033[0m---%p---%T---%+v\n", obj, obj, obj)
	}
}

func DebugPrintV2(module string, obj ...interface{}) {
	fmt.Printf("\033[1;40;32m[debug]\033[0m------[%s] start\n", module)
	for _, v := range obj {
		fmt.Printf("\033[1;40;32m[debug]\033[0m---%p---%T---%+v\n", v, v, v)
	}
	fmt.Printf("\033[1;40;32m[debug]\033[0m------[%s] end\n", module)
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

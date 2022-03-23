package debugutil

import (
	"fmt"
	"os"
)

func DebugPrint(obj interface{}, isExit int) {
	if isExit == 1 {
		fmt.Printf("[debug]---%p---%T---%+v\n", obj, obj, obj)
		os.Exit(0)
	} else {
		fmt.Printf("[debug]---%p---%T---%+v\n", obj, obj, obj)
	}
}

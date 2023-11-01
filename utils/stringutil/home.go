package stringutil

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetAppRootPath() string {
	innerFunc := func(cmdDepth, callerDepth int) string {
		var (
			filePath     string
			defaultDepth int
		)
		binFile, _ := os.Executable()
		//解决不同架构下软连问题，比如mac下会有问题，linux下没有
		binFile, _ = filepath.EvalSymlinks(binFile)
		filePath = path.Dir(binFile)
		defaultDepth = cmdDepth
		if strings.Contains(filePath, "/go-build") || strings.Contains(filePath, "/tmp/GoLand") {
			_, file, _, _ := runtime.Caller(0)
			filePath = path.Dir(file)
			defaultDepth = callerDepth
		}
		for i := 0; i < defaultDepth; i++ {
			filePath += "/.."
		}
		filePath += "/"
		return path.Dir(filePath)
	}
	return innerFunc(1, 3)
}

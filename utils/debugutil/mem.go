package debugutil

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"runtime"
)

func PrintMemStatHuman(m *runtime.MemStats) string {
	if m == nil {
		m = &runtime.MemStats{}
		runtime.ReadMemStats(m)
	} else {
		runtime.ReadMemStats(m)
	}
	ret := fmt.Sprintf(
		"\n堆对象分配大小(Alloc|HeapAlloc): %s\n"+
			"堆对象分配数量(HeapObjects): %d\n"+
			"堆对象累计分配大小(TotalAlloc): %s\n"+
			"堆对象累计分配次数(Mallocs): %d\n"+
			"堆对象累计回收次数(Frees): %d\n",
		humanize.IBytes(m.Alloc),
		m.HeapObjects,
		humanize.IBytes(m.TotalAlloc),
		m.Mallocs,
		m.Frees,
	)
	return ret
}

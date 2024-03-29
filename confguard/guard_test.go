package confguard

import (
	"fmt"
	"github.com/MisterChing/go-lib/utils/debugutil"
	"github.com/MisterChing/go-lib/utils/workflow"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type DemoConf struct {
	Name string
	Age  int
}

func TestNewNacosWatch(t *testing.T) {
	addr := "127.0.0.1:8848"
	namespaceId := "7e26e4a2-a33c-4042-ba2b-dcd63fa05e46"

	var ccc DemoConf
	guarder := NewGuarder(&ccc)

	guard, err := NewGuard(addr, namespaceId, LogLevelDebug,
		WithGroup("dev"),
		WithDataID("service_route_rule.json"),
		WithGuarder(guarder),
		WithWatchKey("serviceRoutes"),
		WithLogger(log.DefaultLogger),
	)
	if err != nil {
		panic(err)
	}
	if err := guard.Watch(); err != nil {
		panic(err)
	}
	defer guard.Close()

	debugutil.DebugPrintV2("outer", guarder, guarder.Get().(*DemoConf))

}

func TestPopulate(t *testing.T) {
	addr := "127.0.0.1:8848"
	namespaceId := "7e26e4a2-a33c-4042-ba2b-dcd63fa05e46"

	var ccc DemoConf
	guarder := NewGuarder(ccc)

	guard, err := NewGuard(addr, namespaceId, LogLevelDebug,
		WithGroup("dev"),
		WithDataID("ddd.json"),
		WithGuarder(guarder),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		jsonStrTpl := `{"name":"ching","age":%d}`
		index := 1
		for {
			jsonStr := fmt.Sprintf(jsonStrTpl, index)
			_ = guard.populate([]byte(jsonStr))
			index++
		}
	}()
	time.Sleep(time.Second)
	index := 1
	for {
		if index > 10000 {
			return
		}
		aa := guarder.Get().(DemoConf)
		debugutil.DebugPrint(aa.Age, 0)
		index++
	}

}

func TestPopulateRace(t *testing.T) {
	addr := "127.0.0.1:8848"
	namespaceId := "7e26e4a2-a33c-4042-ba2b-dcd63fa05e46"

	var ccc DemoConf
	guarder := NewGuarder(ccc)

	guard, err := NewGuard(addr, namespaceId, LogLevelDebug,
		WithGroup("dev"),
		WithDataID("ddd.json"),
		WithGuarder(guarder),
	)
	if err != nil {
		panic(err)
	}

	var count int64
	var wg sync.WaitGroup
	workflow.GoWithNumber(&wg, 10000, func() {
		atomic.AddInt64(&count, 1)
		jsonStrTpl := `{"name":"ching","age":%d}`
		jsonStr := fmt.Sprintf(jsonStrTpl, atomic.LoadInt64(&count))
		_ = guard.populate([]byte(jsonStr))

		aa := guarder.Get().(DemoConf)
		debugutil.DebugPrint(aa.Age, 0)

	})

	debugutil.DebugPrintV2("ret", count, guarder.Get())

}

func TestPopulatePtr(t *testing.T) {
	var ccc DemoConf
	guarder := NewGuarder(&ccc)

	debugutil.DebugPrintV2("00000", guarder, guarder.Get())
	jsonStrTpl := `{"name":"ching","age":%d}`
	jsonStr := fmt.Sprintf(jsonStrTpl, 1)
	_ = guarder.populateWhenPtr([]byte(jsonStr))
	debugutil.DebugPrintV2("11111", guarder, guarder.Get())

}

func TestPopulatePtrLoop(t *testing.T) {
	var ccc DemoConf
	guarder := NewGuarder(&ccc)
	go func() {
		jsonStrTpl := `{"name":"ching","age":%d}`
		index := 1
		for {
			jsonStr := fmt.Sprintf(jsonStrTpl, index)
			//guarder.Lock()
			_ = guarder.populateWhenPtr([]byte(jsonStr))
			//guarder.Unlock()
			index++
		}
	}()
	time.Sleep(time.Second)
	index := 1
	//aa := guarder.Get().(*DemoConf)

	for {
		if index > 100000 {
			return
		}
		aa := guarder.Get().(*DemoConf)

		debugutil.DebugPrint(aa.Age, 0)
		index++
	}

}

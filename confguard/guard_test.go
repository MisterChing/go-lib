package confguard

import (
	"fmt"
	"github.com/MisterChing/go-lib/utils/debugutil"
	"github.com/go-kratos/kratos/v2/log"
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
		if index > 1000000 {
			return
		}
		aa := guarder.Get().(DemoConf)
		debugutil.DebugPrint(aa.Age, 0)
		index++
	}

}

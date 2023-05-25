package confguard

import (
	"encoding/json"
	"fmt"
	knacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	kencoding "github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"reflect"
	"time"
	"unsafe"
)

type Guard struct {
	opts        options
	conf        config.Config
	logger      log.Logger
	logH        *log.Helper
	nacosClient config_client.IConfigClient
}

func NewGuard(addr, namespaceId string, loglevel LogLevel, opts ...Option) (*Guard, error) {
	var (
		nacosSourceArr []config.Source
	)
	_options := options{}
	for _, o := range opts {
		o(&_options)
	}
	nacosClient, err := NewNacosClient(addr, namespaceId, loglevel)
	if err != nil {
		return nil, err
	}
	obj := &Guard{
		opts:        _options,
		nacosClient: nacosClient,
		logger:      _options.logger,
		logH:        log.NewHelper(log.With(_options.logger, "x_module", "conf_guard/conf_guard")),
	}
	tmpSource := knacos.NewConfigSource(nacosClient, knacos.WithGroup(_options.group), knacos.WithDataID(_options.dataID))
	nacosSourceArr = append(nacosSourceArr, tmpSource)

	c := config.New(
		config.WithSource(nacosSourceArr...),
	)
	obj.conf = c
	return obj, nil
}

func (g *Guard) Watch() error {
	if err := g.conf.Load(); err != nil {
		err = fmt.Errorf("confguard load config failed err:%+v", err)
		g.logH.Error(err)
		return err
	}
	var item interface{}
	if err := g.conf.Scan(&item); err != nil {
		err = fmt.Errorf("confguard config unmarshal failed err:%+v", err)
		g.logH.Error(err)
		return err
	} else {
		inputByte, _ := json.Marshal(item)
		if e := g.populate(inputByte); e != nil {
			err = fmt.Errorf("confguard populate failed err:%+v", e)
			g.logH.Error(err)
			return err
		}
	}
	safeWatch(func() error {
		err := g.conf.Watch(g.opts.watchKey, func(key string, value config.Value) {
			var newItem interface{}
			if err := g.conf.Scan(&newItem); err != nil {
				err = fmt.Errorf("confguard config unmarshal failed err:%+v", err)
				g.logH.Error(err)
			} else {
				inputByte, _ := json.Marshal(newItem)
				g.logH.Infof("confguard key:%s changed value:%s", g.opts.watchKey, string(inputByte))

				if e := g.populate(inputByte); e != nil {
					err = fmt.Errorf("confguard populate failed err:%+v", e)
					g.logH.Error(err)
				}
			}
		})
		if err != nil {
			err = fmt.Errorf("confguard watch [key==>%s] config error err:%+v", g.opts.watchKey, err)
			g.logH.Warn(err)
			return err
		}
		return nil
	})
	return nil
}

func (g *Guard) Close() {
	_ = g.conf.Close()
}

func (g *Guard) populate(input []byte) error {
	rv := reflect.ValueOf(g.opts.guarder)
	//g.opts.guarder.lock()
	//defer g.opts.guarder.unlock()
	rv.MethodByName("Lock").Call(nil)         //lock
	defer rv.MethodByName("Unlock").Call(nil) //unlock

	//根据child类型创建一个copy的ptr
	addressableChildCopy := reflect.New(rv.Elem().FieldByName("child").Elem().Type())

	if err := kencoding.GetCodec("json").Unmarshal(input, addressableChildCopy.Interface()); err != nil {
		return err
	}

	//rv.Elem().FieldByName("child").Set(addressableChildCopy.Elem()) //可导出字段可用，非导出字段panic
	//获取非导出字段并转换为可寻址的ptr
	childAddressable := rv.Elem().FieldByName("child")
	childAddressablePtr := reflect.NewAt(childAddressable.Type(), unsafe.Pointer(childAddressable.UnsafeAddr()))
	childAddressablePtr.Elem().Set(addressableChildCopy.Elem())

	return nil
}

func safeWatch(fn func() error) {
	go func() {
		for err := fn(); err != nil; err = fn() {
			time.Sleep(time.Second)
		}
	}()
}

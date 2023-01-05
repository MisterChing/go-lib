package confguard

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
	"strings"
)

func NewNacosClient(addr, namespaceId string, loglevel LogLevel) (config_client.IConfigClient, error) {
	addrSlice := strings.Split(addr, ",")
	sc := make([]constant.ServerConfig, 0, len(addr))
	for _, v := range addrSlice {
		item := strings.Split(v, ":")
		sc = append(sc,
			*constant.NewServerConfig(item[0], cast.ToUint64(item[1])),
		)
	}
	cc := &constant.ClientConfig{
		NamespaceId:         namespaceId, //namespace id
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            string(loglevel),
	}
	// a more graceful way to create naming client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	return client, err
}

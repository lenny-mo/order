package conf

import (
	"fmt"
	"strconv"
	"time"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-plugins/config/source/consul/v2"
)

func GetConfig(host string, port int64, prefix string) (config.Config, error) {

	//
	consulConf := consul.NewSource(
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)), // 这是一个配置选项，用于指定 Consul 服务器的地址。
		consul.WithPrefix(prefix),                                // 用于指定配置的前缀。Consul 配置中心通常使用前缀来组织配置项，本项目使用 /micro/config 作为前缀。
		consul.StripPrefix(true),                                 // 设置为 true 意味着从 Consul 获取配置项时，将从键中去除指定的前缀，好处是直接通过这些简化的键名来访问配置值，而不必处理完整的键名
		// 例如：访问 "host" 而不是 "micro/config/mysql.json/host"
	)

	// 初始化 conf
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	// 加载配置
	err = config.Load(consulConf)
	fmt.Println("加载配置")
	fmt.Println(config)

	time.Sleep(time.Second * 1)
	return config, err
}

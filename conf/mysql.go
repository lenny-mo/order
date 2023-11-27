package conf

import "github.com/micro/go-micro/v2/config"

type MysqlConfig struct {
	Host     string `"json:host" yaml:"host"`
	Port     int64  `"json:port" yaml:"port"`
	User     string `"json:user" yaml:"user"`
	Password string `"json:password" yaml:"password"`
	DB       string `"json:db" yaml:"db"`
}

// GetMysqlFromConsul 从 Consul 配置中心获取 MySQL 配置。
// 它接受一个 config.Config 类型的参数和一个可变长的字符串数组。
// config.Config 通常是一个用于与 Consul 交互的配置对象。
// 可变长的字符串数组 path 用于指定在 Consul 中查找 MySQL 配置的路径。
func GetMysqlFromConsul(config config.Config, path ...string) *MysqlConfig {
	// 创建一个 MysqlConfig 类型的指针 mysqlConfig，用于存储从 Consul 获取的配置。
	mysqlConfig := &MysqlConfig{}

	// 使用 config.Get 方法获取指定路径的配置。
	// path... 是一个语法糖，它将 path 切片展开为一个个独立的参数。
	// Scan 方法将获取到的配置映射（或解析）到 mysqlConfig 对象中。
	config.Get(path...).Scan(mysqlConfig)

	// 返回填充好的 mysqlConfig 对象。
	return mysqlConfig
}

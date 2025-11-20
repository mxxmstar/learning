package config

type Config struct {
	Env string `mapstructure:"ENV"`
	// 添加其它需要的配置项
}

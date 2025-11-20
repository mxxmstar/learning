package config

import (
	rootconfig "github.com/mxxmstar/learning.git/config"
)

// Config 类型别名
type Config = rootconfig.Config

func Init() (*Config, error) {
	return rootconfig.Init()
}

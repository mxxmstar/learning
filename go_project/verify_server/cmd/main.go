package main

import (
	"fmt"
	"log"

	"github.com/mxxmstar/learning/verify_server/internal/web"
	"github.com/mxxmstar/learning/verify_server/verify_config"
)

func main() {
	fmt.Println("Hello, World!")
	// 初始化配置
	cfg, err := verify_config.Init()
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}

	// 打印配置
	log.Printf("Config: %+v\n", cfg)
	// 初始化服务
	server := web.InitWebServer(cfg)
	// 启动服务
	// server.Run(fmt.Sprintf(":%d", cfg.VerifyServer.Port))
	server.Run(fmt.Sprintf("0.0.0.0:%d", cfg.ServerConfig.VerifyServers.Port))
}

package main

import (
	"fmt"
	"log"

	"github.com/mxxmstar/learning/verify_server/config"
)

func main() {
	fmt.Println("Hello, World!")
	// 初始化配置
	cfg, err := config.Init()
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		return
	}

	// 打印配置
	log.Printf("Config: %+v\n", cfg)
}

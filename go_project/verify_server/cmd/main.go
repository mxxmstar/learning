package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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

func initWebServer(cfg *config.Config) error {
	server := gin.Default()
	server.Use(func(c *gin.Context) {

	})

}

package main

import "github.com/mxxmstar/learning/status_server/status_config"

func main() {
	cfg, err := status_config.Init()
	if err != nil {
		panic(err)
	}

}

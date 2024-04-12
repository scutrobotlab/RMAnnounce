package main

import "git.scutbot.cn/Web/RMAnnounce/internal/config"

func main() {
	c := config.NewConfig("etc/config.yaml")
	if c == nil {
		panic("Failed to load config")
	}
}

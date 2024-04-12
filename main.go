package main

import (
	"git.scutbot.cn/Web/RMAnnounce/internal/config"
	"git.scutbot.cn/Web/RMAnnounce/internal/job"
)

func main() {
	c := config.NewConfig()
	if c == nil {
		_ = config.Config{}.Save()
		panic("Failed to load config")
	}

	cron := job.InitCronjob()
	if cron == nil {
		panic("Failed to init cronjob")
	}
	cron.Start()
	defer cron.Stop()

	select {}
}

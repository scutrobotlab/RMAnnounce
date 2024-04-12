package main

import (
	"git.scutbot.cn/Web/RMAnnounce/internal/config"
	"git.scutbot.cn/Web/RMAnnounce/internal/job"
)

func main() {
	c := config.NewConfig("etc/config.yaml")
	if c == nil {
		panic("Failed to load config")
	}
	config.Instance = *c

	cron := job.InitCronjob()
	if cron == nil {
		panic("Failed to init cronjob")
	}
	cron.Start()
	defer cron.Stop()

	select {}
}

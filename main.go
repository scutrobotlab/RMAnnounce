package main

import (
	"fmt"
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"github.com/scutrobotlab/RMAnnounce/internal/job"
)

func main() {
	c := config.NewConfig(config.DefaultPath)
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

	fmt.Printf("RoboMaster Announce Bot started\n")

	select {}
}

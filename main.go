package main

import (
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"github.com/scutrobotlab/RMAnnounce/internal/job"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.NewConfig(config.DefaultPath)
	if c == nil {
		_ = config.Config{}.Save()
		c = &config.Config{}
		logrus.Infof("Config file not found, a new one has been created")
	}

	cron := job.InitCronjob()
	if cron == nil {
		panic("Failed to init cronjob")
	}
	cron.Start()
	defer cron.Stop()

	logrus.Infof("RoboMaster Announce Bot started")

	select {}
}

package job

import (
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

// InitCronjob initializes the cronjob
func InitCronjob() *cron.Cron {
	c := cron.New()

	fetchAnnounce := &FetchAnnounceJob{
		SentMap: cache.New(cache.DefaultExpiration, cache.DefaultExpiration),
	}
	_, err := c.AddJob("@every 15s", fetchAnnounce)
	fetchAnnounce.Init()
	if err != nil {
		logrus.Errorf("Failed to add job: %v", err)
		return nil
	}

	monitorAnnounce := MonitorAnnounceJob{}
	_, err = c.AddJob("@every 15s", monitorAnnounce)
	monitorAnnounce.Init()
	if err != nil {
		logrus.Errorf("Failed to add job: %v", err)
		return nil
	}

	return c
}

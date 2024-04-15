package job

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

// InitCronjob initializes the cronjob
func InitCronjob() *cron.Cron {
	c := cron.New()

	fetchAnnounce := FetchAnnounceJob{}
	_, err := c.AddJob("@every 5s", fetchAnnounce)
	fetchAnnounce.Init()
	if err != nil {
		fmt.Printf("Failed to add job: %v\n", err)
		return nil
	}

	return c
}

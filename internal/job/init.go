package job

import (
	"fmt"
	"git.scutbot.cn/Web/RMAnnounce/internal/config"
	"github.com/robfig/cron/v3"
)

// InitCronjob initializes the cronjob
func InitCronjob() *cron.Cron {
	c := cron.New()
	_, err := c.AddJob("@every 5s", FetchAnnounceJob{c: config.GetInstance()})
	if err != nil {
		fmt.Printf("Failed to add job: %v\n", err)
		return nil
	}

	return c
}

package job

import "git.scutbot.cn/Web/RMAnnounce/internal/config"

type FetchAnnounceJob struct {
	c config.Config
}

func (f FetchAnnounceJob) Run() {
	println("Fetching announce...")
}

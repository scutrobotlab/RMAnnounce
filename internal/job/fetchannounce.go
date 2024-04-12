package job

import (
	"fmt"
	"git.scutbot.cn/Web/RMAnnounce/internal/config"
	"git.scutbot.cn/Web/RMAnnounce/internal/util"
	"io"
	"log"
	"net/http"
	"strings"
)

type FetchAnnounceJob struct {
}

func (f FetchAnnounceJob) Run() {
	c := config.GetInstance()
	url := getUrl(c.LastId + 1)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	bodyStr := string(body)
	if strings.Contains(bodyStr, "您访问的页面不存在") {
		log.Printf("Page %d not found", c.LastId+1)
		return
	}

	log.Printf("Found new announcement: %d", c.LastId+1)
	c.LastId++
	err = c.Save()
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := fmt.Sprintf("RoboMaster 资料站新公告\n%s", url)
	err = util.SendWebhookMsg(c.Webhook, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getUrl(id int) string {
	return fmt.Sprintf("https://www.robomaster.com/zh-CN/resource/pages/announcement/%d", id)
}

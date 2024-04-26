package job

import (
	"fmt"
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"github.com/scutrobotlab/RMAnnounce/internal/util"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

type FetchAnnounceJob struct {
}

func (f FetchAnnounceJob) Init() {
	c := config.GetInstance()
	fmt.Printf("Load webhooks count: %d\n", len(c.Webhooks))
}

func (f FetchAnnounceJob) Run() {
	c := config.GetInstance()
	if c.LastId == 0 {
		fmt.Printf("LastId is 0, skip\n")
		return
	}

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

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(bodyStr))
	if err != nil {
		log.Fatal(err)
		return
	}
	title, err := getMainTitle(doc)
	if err != nil {
		log.Println(err)
		return
	}

	contents := [][]util.Content{
		{
			{
				Tag:    "at",
				UserId: "all",
			},
			{
				Tag:  "text",
				Text: " " + title + "\n",
			},
			{
				Tag:  "text",
				Text: url,
			},
		},
	}
	err = util.SendPostMsg(c.Webhooks, "RoboMaster 资料站新公告", contents)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getUrl(id int) string {
	return fmt.Sprintf("https://www.robomaster.com/zh-CN/resource/pages/announcement/%d", id)
}

// 递归查找主标题
func getMainTitle(n *html.Node) (string, error) {
	// 如果是p标签且class为main-title，则返回标题
	if n.Type == html.ElementNode && n.Data == "p" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "main-title" {
				return n.FirstChild.Data, nil
			}
		}
	}

	// 递归处理子节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, err := getMainTitle(c)
		if err == nil {
			return title, nil
		}
	}

	return "", fmt.Errorf("main title not found")
}

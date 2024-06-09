package job

import (
	"crypto/sha256"
	"fmt"
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"github.com/scutrobotlab/RMAnnounce/internal/util"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

type MonitorAnnounceJob struct {
}

func (m MonitorAnnounceJob) Init() {
	c := config.GetInstance()
	log.Printf("Monitor pages count: %d\n", len(c.MonitoredPages))
}

func (m MonitorAnnounceJob) Run() {
	c := config.GetInstance()
	for i, page := range c.MonitoredPages {
		url := getUrl(page.Id)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to get page %d: %v", page.Id, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Failed to read page %d: %v", page.Id, err)
			continue
		}

		bodyStr := string(body)
		doc, err := html.Parse(strings.NewReader(bodyStr))
		if err != nil {
			log.Printf("Failed to parse page %d: %v", page.Id, err)
			continue
		}

		mainContext, err := getMainContext(doc)
		if err != nil {
			log.Printf("Failed to get main context of page %d: %v", page.Id, err)
			continue
		}

		hash, err := getMainContextHash(mainContext)
		if err != nil {
			log.Printf("Failed to get hash of page %d: %v", page.Id, err)
			continue
		}

		if page.Hash == "" {
			c.MonitoredPages[i].Hash = hash
			err = c.Save()
			if err != nil {
				log.Printf("Failed to save page %d: %v", page.Id, err)
				continue
			}

			log.Printf("Init hash of page %d: %s", page.Id, hash)
			continue
		}

		if page.Hash != hash {
			c.MonitoredPages[i].Hash = hash
			err = c.Save()
			if err != nil {
				log.Printf("Failed to save page %d: %v", page.Id, err)
				continue
			}

			log.Printf("Hash changed of page %d: %s", page.Id, hash)
			var title string
			title, err = getMainTitle(doc)
			if err != nil {
				log.Println(err)
				continue
			}

			contents := [][]util.Content{
				{
					{
						Tag:    "at",
						UserId: "all",
					},
					{
						Tag:  "text",
						Text: " [更新] " + title + "\n",
					},
					{
						Tag:  "text",
						Text: url,
					},
				},
			}
			err = util.SendPostMsg(c.Webhooks, "RoboMaster 资料站新公告", contents)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

// 递归查找主内容
func getMainContext(n *html.Node) (*html.Node, error) {
	// 如果是p标签且class为main-title，则返回标题
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "main-context" {
				return n, nil
			}
		}
	}

	// 递归处理子节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		context, err := getMainContext(c)
		if err == nil {
			return context, nil
		}
	}

	return nil, fmt.Errorf("main context not found")
}

// 获取主内容的 HASH
func getMainContextHash(node *html.Node) (string, error) {
	var w strings.Builder
	err := html.Render(&w, node)
	if err != nil {
		return "", err
	}

	str := w.String()
	sum := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%X", sum), nil
}

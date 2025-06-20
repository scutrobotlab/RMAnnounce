package job

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"github.com/scutrobotlab/RMAnnounce/internal/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
	"time"
)

type FetchAnnounceJob struct {
	SentMap *cache.Cache // 用于存储已发送的公告ID，避免重复发送
}

func (f *FetchAnnounceJob) Init() {
	c := config.GetInstance()
	logrus.Infof("Load webhooks count: %d", len(c.Webhooks))
}

func (f *FetchAnnounceJob) Run() {
	c := config.GetInstance()
	if c.LastId == 0 {
		logrus.Infof("LastId is 0, skip")
		return
	}

	nextId := c.LastId + 1
	nextIdStr := fmt.Sprintf("%d", nextId)
	_, ok := f.SentMap.Get(nextIdStr)
	if ok {
		logrus.Warnf("Announcement %d already sent, skipping", nextId)
		c.LastId = nextId
		err := c.Save()
		if err != nil {
			logrus.Errorf("Failed to save config: %v", err)
			return
		}
		logrus.Infof("Updated LastId to %d", c.LastId)
		return
	}
	url := getUrl(nextId)

	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("Failed to http get page %d: %v", nextId, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			logrus.Infof("Page %d not found because of 404", nextId)
			return
		}
		// 其他错误状态码
		logrus.Errorf("Failed to fetch page %d: status code %d", nextId, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read page %d: %v", nextId, err)
		return
	}

	bodyStr := string(body)
	if strings.Contains(bodyStr, "您访问的页面不存在") {
		logrus.Infof("Page %d not found because of '您访问的页面不存在'", nextId)
		return
	}

	logrus.Infof("Found new announcement: %d", nextId)
	c.LastId++
	err = c.Save()
	if err != nil {
		logrus.Errorf("Failed to save config: %v", err)
		return
	}

	// 解析HTML
	doc, err := html.Parse(strings.NewReader(bodyStr))
	if err != nil {
		logrus.Errorf("Failed to parse page %d: %v", nextId, err)
		return
	}
	title, err := getMainTitle(doc)
	if err != nil {
		logrus.Errorf("Failed to get main title of page %d: %v", nextId, err)
		return
	}

	mainContext, err := getMainContext(doc)
	if err != nil {
		logrus.Errorf("Failed to get main context of page %d: %v", nextId, err)
		return
	}
	contextIsEmpty := mainContext.FirstChild == nil

	var atAllType util.AtAllStatus
	var contents [][]util.Content
	if contextIsEmpty {
		atAllType = util.AtAllStatusFalse
		contents = [][]util.Content{
			{
				{
					Tag:  "text",
					Text: "[空白] " + title + "\n",
				},
				{
					Tag:  "text",
					Text: url,
				},
			},
		}
	} else {
		atAllType = util.AtAllStatusAuto
		contents = [][]util.Content{
			{
				{
					Tag:  "text",
					Text: "[新增] " + title + "\n",
				},
				{
					Tag:  "text",
					Text: url,
				},
			},
		}
	}

	ok, err = util.SendPostMsg(c.Webhooks, "RoboMaster 资料站新公告", atAllType, contents)
	if err != nil {
		logrus.Errorf("Failed to send robotomaster notification: %v", err)
		return
	}
	if ok {
		logrus.Infof("Announcement %d sent successfully", nextId)
		f.SentMap.Set(nextIdStr, struct{}{}, time.Hour)
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

package job

import (
	"crypto/sha256"
	"fmt"
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"golang.org/x/net/html"
	"strings"
)

type MonitorAnnounceJob struct {
}

func (m MonitorAnnounceJob) Init() {
	c := config.GetInstance()
	fmt.Printf("Monitor pages count: %d\n", len(c.MonitoredPages))
}

func (m MonitorAnnounceJob) Run() {
	fmt.Printf("MonitorAnnounceJob run\n")
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

package job

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestGetMainTitle(t *testing.T) {
	// 获取HTML内容
	resp, err := http.Get("https://www.robomaster.com/zh-CN/resource/pages/announcement/1708")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 解析HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 从根节点开始查找
	title, err := getMainTitle(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Main title: %s\n", title)
}

func TestGetMainContext(t *testing.T) {
	// 获取HTML内容
	resp, err := http.Get("https://www.robomaster.com/zh-CN/resource/pages/announcement/1653")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 解析HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 从根节点开始查找
	context, err := getMainContext(doc)
	if err != nil {
		log.Fatal(err)
	}

	var w strings.Builder
	err = html.Render(&w, context)
	if err != nil {
		log.Fatal(err)
	}

	str := w.String()
	fmt.Printf("Main context: %s\n", str)
}

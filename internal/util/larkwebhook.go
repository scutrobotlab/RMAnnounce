package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
	"time"
)

type SendCircuitBreaker struct {
	Name     string
	Counter  *SlidingWindowCounter
	MaxCount int
}

var (
	SendCircuitBreakers = []SendCircuitBreaker{
		{
			Name:     "每小时",
			Counter:  NewSlidingWindowCounter(time.Hour),
			MaxCount: 15,
		},
		{
			Name:     "每分钟",
			Counter:  NewSlidingWindowCounter(time.Minute),
			MaxCount: 5,
		},
		{
			Name:     "每5秒",
			Counter:  NewSlidingWindowCounter(time.Second * 5),
			MaxCount: 3,
		},
	}
	CircuitBreakerRecoverTime = atomic.Int64{}
)

const CircuitBreakerDuration = time.Hour * 12 // 断路器恢复时间

// LastAtAllTime 最后一次@所有人的时间戳
var LastAtAllTime atomic.Int64

// AtAllAutoInterval 自动@所有人的时间间隔
const AtAllAutoInterval = time.Second * 60

type AtAllStatus int

const (
	AtAllStatusFalse AtAllStatus = iota // 不@所有人
	AtAllStatusTrue  AtAllStatus = iota // @所有人
	AtAllStatusAuto  AtAllStatus = iota // 自动判断
)

type WebhookBotTextReq struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type WebhookBotResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type WebhookBotPostReq struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Post struct {
			ZhCn struct {
				Title   string      `json:"title"`
				Content [][]Content `json:"content"`
			} `json:"zh_cn"`
		} `json:"post"`
	} `json:"content"`
}

type Content struct {
	Tag    string `json:"tag"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

// SendTextMsg 发送文本消息
func SendTextMsg(urls []string, msg string) error {
	req := WebhookBotTextReq{
		MsgType: "text",
	}
	req.Content.Text = msg

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return SendWebhookMsg(urls, reqBody)
}

func CheckCircuitBreaker(urls []string) bool {
	circuitBreakerRecoverTime := time.UnixMilli(CircuitBreakerRecoverTime.Load())
	if time.Now().Before(circuitBreakerRecoverTime) {
		// 断路器处于开启状态，拒绝发送消息
		return true
	}
	for _, breaker := range SendCircuitBreakers {
		breaker.Counter.Increment()
		if breaker.Counter.Count() > breaker.MaxCount {
			// 断路器触发，拒绝发送消息
			logrus.Warnf("Circuit breaker triggered: %s, count: %d, max count: %d", breaker.Name, breaker.Counter.Count(), breaker.MaxCount)
			recoverTime := time.Now().Add(CircuitBreakerDuration)
			CircuitBreakerRecoverTime.Store(recoverTime.UnixMilli())
			// 发送求救消息
			content := [][]Content{
				{
					{Tag: "at", UserId: "ou_892a4b47aa876f799ca3aef97403e009"},
					{Tag: "text", Text: " 已自动触发熔断，请救救我"},
				},
				{
					{Tag: "text", Text: fmt.Sprintf("%s尝试发送数量达到 %d，超过最大限制 %d，已触发熔断。", breaker.Name, breaker.Counter.Count(), breaker.MaxCount)},
					{Tag: "text", Text: fmt.Sprintf("为避免打扰，我在 %s 前会保持安静。", recoverTime.Format(time.DateTime))},
				},
			}
			_, err := sendPostMsg(urls, "机器人碎碎念", content)
			if err != nil {
				logrus.Errorf("send post msg error: %s", err)
			}
			return true
		}
	}
	// 断路器处于关闭状态，允许发送消息
	return false
}

// SendPostMsg 发送富文本消息
func SendPostMsg(urls []string, title string, atAllStatus AtAllStatus, content [][]Content) (bool, error) {
	if CheckCircuitBreaker(urls) {
		logrus.Warnf("Circuit breaker is open, cannot send message until %s", time.UnixMilli(CircuitBreakerRecoverTime.Load()).Format(time.RFC3339))
		return false, nil
	}

	var atAll bool
	switch atAllStatus {
	case AtAllStatusFalse:
		atAll = false
	case AtAllStatusTrue:
		atAll = true
	case AtAllStatusAuto:
		lastAtAllTime := time.UnixMilli(LastAtAllTime.Load())
		atAll = time.Since(lastAtAllTime) >= AtAllAutoInterval
	default:
		return false, errors.New("invalid atAllStatus")
	}
	if atAll {
		if len(content) != 0 {
			content[0] = append([]Content{
				{Tag: "at", UserId: "all"},
				{Tag: "text", Text: " "},
			}, content[0]...)
		} else {
			return false, errors.New("content cannot be empty when atAll is true")
		}
		LastAtAllTime.Store(time.Now().UnixMilli())
	}

	return sendPostMsg(urls, title, content)
}

func sendPostMsg(urls []string, title string, content [][]Content) (bool, error) {
	req := WebhookBotPostReq{
		MsgType: "post",
	}
	req.Content.Post.ZhCn.Title = title
	req.Content.Post.ZhCn.Content = content

	reqBody, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	return true, SendWebhookMsg(urls, reqBody)
}

func SendWebhookMsg(urls []string, body []byte) error {
	for _, url := range urls {
		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			logrus.Errorf("Failed to send webhook msg: %v", err)
			continue
		}

		var respBody WebhookBotResp
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			logrus.Errorf("Failed to send webhook msg: %v", err)
			continue
		}

		if respBody.Code != 0 {
			logrus.Errorf("Failed to send webhook msg: %v", respBody.Msg)
			continue
		}

		resp.Body.Close()
	}

	return nil
}

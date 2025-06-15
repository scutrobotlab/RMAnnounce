package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

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

// SendPostMsg 发送富文本消息
func SendPostMsg(urls []string, title string, atAllStatus AtAllStatus, content [][]Content) error {
	req := WebhookBotPostReq{
		MsgType: "post",
	}
	req.Content.Post.ZhCn.Title = title

	var atAll bool
	switch atAllStatus {
	case AtAllStatusFalse:
		atAll = false
	case AtAllStatusTrue:
		atAll = true
	case AtAllStatusAuto:
		// TODO: 自动判断
	default:
		return errors.New("invalid atAllStatus")
	}
	if atAll {
		if len(content) != 0 {
			content[0] = append([]Content{
				{Tag: "at", UserId: "all"},
				{Tag: "text", Text: " "},
			}, content[0]...)
		} else {
			return errors.New("content cannot be empty when atAll is true")
		}
	}
	req.Content.Post.ZhCn.Content = content

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return SendWebhookMsg(urls, reqBody)
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

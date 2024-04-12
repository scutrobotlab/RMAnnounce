package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type WebhookBotReq struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type WebhookBotResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SendWebhookMsg(urls []string, msg string) error {
	req := WebhookBotReq{
		MsgType: "text",
	}
	req.Content.Text = msg

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	for _, url := range urls {
		resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
		if err != nil {
			log.Printf("Failed to send webhook message: %v", err)
			continue
		}

		var respBody WebhookBotResp
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			log.Printf("Failed to send webhook message: %v", err)
			continue
		}

		if respBody.Code != 0 {
			log.Printf("Failed to send webhook message: %s", respBody.Msg)
			continue
		}

		resp.Body.Close()
	}

	return nil
}

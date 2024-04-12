package util

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func SendWebhookMsg(url string, msg string) error {
	req := WebhookBotReq{
		MsgType: "text",
	}
	req.Content.Text = msg

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respBody WebhookBotResp
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	if respBody.Code != 0 {
		return fmt.Errorf("failed to send webhook message: %s", respBody.Msg)
	}

	return nil
}

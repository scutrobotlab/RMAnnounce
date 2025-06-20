package util

import (
	"github.com/scutrobotlab/RMAnnounce/internal/config"
	"testing"
)

func TestSendWebhookMsg(t *testing.T) {
	c := config.NewConfig("../../etc/config.yaml")
	if c == nil {
		t.Fatal("Failed to load config")
	}

	_, err := SendPostMsg(c.Webhooks, "机器人碎碎念", AtAllStatusFalse, [][]Content{
		{
			//{Tag: "at", UserId: "ou_892a4b47aa876f799ca3aef97403e009"},
			{Tag: "text", Text: "测试富文本"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

func TestSendWebhookMsgAtAll(t *testing.T) {
	c := config.NewConfig("../../etc/config.yaml")
	if c == nil {
		t.Fatal("Failed to load config")
	}

	_, err := SendPostMsg(c.Webhooks, "机器人碎碎念", AtAllStatusTrue, [][]Content{
		{
			{Tag: "text", Text: "测试富文本"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

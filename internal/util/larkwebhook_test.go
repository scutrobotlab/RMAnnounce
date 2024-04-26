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

	err := SendPostMsg(c.Webhooks, "机器人碎碎念", [][]Content{
		{
			{Tag: "at", UserId: "ou_892a4b47aa876f799ca3aef97403e009"},
			{Tag: "text", Text: " TEXT"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

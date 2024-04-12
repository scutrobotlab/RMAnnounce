package util

import (
	"git.scutbot.cn/Web/RMAnnounce/internal/config"
	"testing"
)

func TestSendWebhookMsg(t *testing.T) {
	c := config.NewConfig("../../etc/config.yaml")
	if c == nil {
		t.Fatal("Failed to load config")
	}

	err := SendWebhookMsg(c.Webhook, "test")
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

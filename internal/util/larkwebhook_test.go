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

	err := SendWebhookMsg(c.Webhooks, "RoboMaster 资料站新公告\n大家好，我是RM资料站公告推送机器人。完整形态考核成绩发布时，我会第一时间通知大家，敬请期待！")
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

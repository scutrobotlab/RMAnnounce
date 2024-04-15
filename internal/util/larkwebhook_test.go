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

	err := SendWebhookMsg(c.Webhooks, "RoboMaster 资料站新公告\n"+
		"(¦3[▓▓] 趁主人不注意，我偷偷睡着了。这次没有第一时间更新公告，抱歉抱歉！(。・＿・。)ﾉI’m sorry~\n")
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

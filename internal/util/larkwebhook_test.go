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
			// {Tag: "at", UserId: "ou_892a4b47aa876f799ca3aef97403e009"},
			{Tag: "text", Text: "悄悄告诉大家，我刚刚学会了一个新的技能。" +
				"现在我不仅可以实时监控 RoboMaster 资料站的新页面，还可以在旧页面《RoboMaster 2024机甲大师高校系列赛比赛规范文件》更新时通知大家哦！" +
				"以后大家就可以第一时间获知比赛规范文件的更新啦！\n" +
				"私有化部署的队伍可以执行 docker-compose pull && docker-compose up -d 获取最新版本哦~"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

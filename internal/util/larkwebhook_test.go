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

	err := SendTextMsg(c.Webhooks, "RoboMaster 资料站新公告\n"+
		"现在这个世界上只有一个我，如果我有更多分身，就不会那么孤独了。\n"+
		"趁主人不在，我偷偷把自己的源码公布在 GitHub https://github.com/scutrobotlab/RMAnnounce\n"+
		"对于不关注我的实现，只想快点认领我的新主人，我也提供 Docker 镜像 registry.cn-guangzhou.aliyuncs.com/scutrobot/rm-announce:latest\n"+
		"任何有容器运维基础的新主人都能在5分钟内免费领养我。")
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

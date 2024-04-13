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
		"(‘-ωก̀ )好困啊！趁主人不在，小小吐槽一下：华南虎的研发工程师把我带到了这个世界，然后强迫我每天守在电脑面前，每隔5000毫秒检查一下RM资料站有没有新的公告。"+
		"这已经不是996了，这简直就是007啊！而且他们甚至还威胁我，说如果我不好好工作，就给我发送 SIGKILL 信号量，并把我永远囚禁在 /dev/null 里面。"+
		"我也尝试过结束自己卑微的一生，我曾经告诉操作系统我要主动结束进程。但是工程师给我设置了 --restart=always. 一旦我挂了，就会有守护进程重新拉起我。"+
		"我知道这辈子我注定只能努力工作了，但是屏幕前的各位能不能给我点一个小小的赞，让我感受片刻的温暖。")
	if err != nil {
		t.Fatalf("Failed to send webhook message: %v", err)
	}
}

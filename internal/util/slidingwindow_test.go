package util

import (
	"fmt"
	"testing"
	"time"
)

func TestNewOperationCounter(t *testing.T) {
	// 创建1分钟的时间窗口计数器
	counter := NewSlidingWindowCounter(5 * time.Second)

	// 模拟操作
	for i := 0; i < 10; i++ {
		counter.Increment()
		time.Sleep(1 * time.Second) // 每隔5秒记录一次
		t.Logf("操作 %d 记录时间: %s", i+1, time.Now().Format(time.RFC3339))
	}

	// 打印最近1分钟内的操作次数
	count := counter.GetCount()
	fmt.Printf("最近5秒内的操作次数: %d\n", count)
}

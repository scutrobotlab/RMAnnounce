package util

import (
	"sync"
	"time"
)

// SlidingWindowCounter 是一个滑动窗口计数器，用于统计在指定时间窗口内的操作次数
type SlidingWindowCounter struct {
	mutex      sync.Mutex    // 互斥锁，保证并发安全
	timestamps []time.Time   // 存储操作时间戳
	window     time.Duration // 时间窗口大小
}

// NewSlidingWindowCounter 创建新的滑动窗口计数器
func NewSlidingWindowCounter(window time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		timestamps: []time.Time{},
		window:     window,
	}
}

// Increment 记录一次操作
func (oc *SlidingWindowCounter) Increment() {
	oc.mutex.Lock()
	defer oc.mutex.Unlock()

	// 添加当前时间戳
	oc.timestamps = append(oc.timestamps, time.Now())
}

// Count 获取时间窗口内的操作次数
func (oc *SlidingWindowCounter) Count() int {
	oc.mutex.Lock()
	defer oc.mutex.Unlock()

	now := time.Now()
	minTime := now.Add(-oc.window)
	count := 0

	// 遍历时间戳，统计在窗口内的次数
	for _, ts := range oc.timestamps {
		if ts.After(minTime) {
			count++
		}
	}

	// 清理过期时间戳（优化内存）
	var cleaned []time.Time
	for _, ts := range oc.timestamps {
		if ts.After(minTime) {
			cleaned = append(cleaned, ts)
		}
	}
	oc.timestamps = cleaned

	return count
}

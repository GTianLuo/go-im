package util

import (
	"github.com/shirou/gopsutil/cpu"
	"time"
)

func CPUPercent() float64 {
	// 获取 CPU 占用率
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}
	return percent[0]
}

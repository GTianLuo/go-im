package util

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"sort"
	"testing"
	"time"
)

func TestCpu(t *testing.T) {
	// 获取 CPU 占用率
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(percent)
}

func TestSort(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	sort.Slice(s, func(i, j int) bool {
		return s[i] > s[j]
	})
	fmt.Println(s)
}

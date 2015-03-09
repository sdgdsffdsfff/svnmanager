package proc

import (
	"math"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func CPUPercent() float64 {
	percent, err := cpu.CPUPercent(time.Second * 2, false)
	if err == nil && len(percent) > 0 {
		return math.Floor(percent[0])
	}
	return 0
}

func MEMPercent() float64{
	v, _ := mem.VirtualMemory()
	return math.Floor(v.UsedPercent)
}


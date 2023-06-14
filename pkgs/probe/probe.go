package probe

import (
	"strconv"

	load "github.com/shirou/gopsutil/v3/load"
	mem "github.com/shirou/gopsutil/v3/mem"
)

func Probe() string {
	res := ""

	load, _ := load.Avg()
	res += "Load Average:\n\t\t\t1m: " + strconv.FormatFloat(load.Load1, 'f', -1, 64) + " 5m: " + strconv.FormatFloat(load.Load5, 'f', -1, 64) + " 15m: " + strconv.FormatFloat(load.Load15, 'f', -1, 64) + "\n\n"

	m, _ := mem.VirtualMemory()
	res += "Memory Usage:\n\t\t\tTotal: " + strconv.FormatUint(m.Total/1024/1024/1024, 10) + "GB Used %: " + strconv.FormatFloat(m.UsedPercent, 'f', 2, 64) + "\n\n"

	return res
}

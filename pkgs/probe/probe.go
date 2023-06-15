package probe

import (
	"log"
	"math"
	"strconv"
	"time"

	disk "github.com/shirou/gopsutil/v3/disk"
	host "github.com/shirou/gopsutil/v3/host"
	load "github.com/shirou/gopsutil/v3/load"
	mem "github.com/shirou/gopsutil/v3/mem"
)

type DisksStats struct {
	Paths []string
}

func parseUptime(uptime uint64) string {
	res := ""
	d, err := time.ParseDuration(strconv.FormatUint(uptime, 10) + "s")
	if err == nil {
		days := d.Hours() / 24
		hours := uint64(d.Hours()) - uint64(days)*24
		hmod := math.Mod(d.Hours(), 60)
		res = strconv.FormatFloat(days, 'f', 0, 64) + " days " + strconv.FormatUint(hours, 10) + " hours " + strconv.FormatFloat(hmod*60/100, 'f', 0, 64) + " minutes"
	} else {
		log.Println(err)
	}
	return res
}

func (data *DisksStats) Probe() string {
	res := ""

	utime, _ := host.Uptime()

	res += "Uptime:\n\t\t\t " + parseUptime(utime) + "\n\n"

	load, _ := load.Avg()
	res += "Load Average:\n\t\t\t1m: " + strconv.FormatFloat(load.Load1, 'f', -1, 64) + " 5m: " + strconv.FormatFloat(load.Load5, 'f', -1, 64) + " 15m: " + strconv.FormatFloat(load.Load15, 'f', -1, 64) + "\n\n"

	m, _ := mem.VirtualMemory()
	res += "Memory Usage:\n\t\t\tUsed: " + strconv.FormatUint(m.Used/1024/1024/1024, 10) + "GB/" + strconv.FormatUint(m.Total/1024/1024/1024, 10) + "GB " + strconv.FormatFloat(m.UsedPercent, 'f', 2, 64) + "%\n\n"

	res += "Disks:\n\t\t\t"
	for _, path := range data.Paths {
		d, _ := disk.Usage(path)
		res += path + ": " + strconv.FormatUint(d.Used/1024/1024/1024, 10) + "GB/" + strconv.FormatUint(d.Total/1024/1024/1024, 10) + "GB " + strconv.FormatFloat(d.UsedPercent, 'f', 2, 64) + "%\n\t\t\t"
	}
	res += "\n"

	if temp, err := host.SensorsTemperatures(); err == nil {
		res += "Temperature:\n\t\t\t"
		for _, t := range temp {
			res += t.SensorKey + ": " + strconv.FormatFloat(t.Temperature, 'f', 2, 64) + " High: " + strconv.FormatFloat(t.High, 'f', 2, 64) + " Critical: " + strconv.FormatFloat(t.Critical, 'f', 2, 64) + "\n\t\t\t"
		}
	}

	return res
}

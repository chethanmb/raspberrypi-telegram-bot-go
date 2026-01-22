package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/* ---------- Temperature ---------- */

func cpuTemp() string {
	out, err := exec.Command("vcgencmd", "measure_temp").Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out)) // temp=48.2'C
}

/* ---------- Uptime ---------- */

func uptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "N/A"
	}

	secs, _ := strconv.ParseFloat(strings.Fields(string(data))[0], 64)
	d := time.Duration(secs) * time.Second

	return fmt.Sprintf(
		"%dd %dh %dm",
		int(d.Hours())/24,
		int(d.Hours())%24,
		int(d.Minutes())%60,
	)
}

/* ---------- CPU Usage ---------- */

func cpuUsage() string {
	idle1, total1 := cpuStat()
	time.Sleep(500 * time.Millisecond)
	idle2, total2 := cpuStat()

	idle := idle2 - idle1
	total := total2 - total1

	usage := 100 * (1 - float64(idle)/float64(total))
	return fmt.Sprintf("%.1f%%", usage)
}

func cpuStat() (idle, total uint64) {
	data, _ := os.ReadFile("/proc/stat")
	fields := strings.Fields(strings.Split(string(data), "\n")[0])[1:]

	for i, v := range fields {
		val, _ := strconv.ParseUint(v, 10, 64)
		total += val
		if i == 3 { // idle
			idle = val
		}
	}
	return
}

/* ---------- RAM Info ---------- */

func ramInfo() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "N/A"
	}

	var total, free uint64

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal") {
			fmt.Sscanf(line, "MemTotal: %d kB", &total)
		}
		if strings.HasPrefix(line, "MemAvailable") {
			fmt.Sscanf(line, "MemAvailable: %d kB", &free)
		}
	}

	used := total - free
	percent := float64(used) / float64(total) * 100

	return fmt.Sprintf(
		"%.1fGB / %.1fGB (%.0f%%)",
		float64(used)/1024/1024,
		float64(total)/1024/1024,
		percent,
	)
}

/* ---------- Throttling ---------- */

func throttled() string {
	out, err := exec.Command("vcgencmd", "get_throttled").Output()
	if err != nil {
		return "N/A"
	}

	valStr := strings.TrimPrefix(strings.TrimSpace(string(out)), "throttled=")
	val, _ := strconv.ParseInt(valStr, 0, 64)

	if val == 0 {
		return "✅OK"
	}

	var issues []string

	if val&0x1 != 0 {
		issues = append(issues, "Under-voltage")
	}
	if val&0x2 != 0 {
		issues = append(issues, "Freq capped")
	}
	if val&0x4 != 0 {
		issues = append(issues, "Throttled")
	}
	if val&0x8 != 0 {
		issues = append(issues, "Soft temp limit")
	}

	return strings.Join(issues, ", ")
}

func cpuTempValue() float64 {
	out, err := exec.Command("vcgencmd", "measure_temp").Output()
	if err != nil {
		return 0
	}

	// temp=48.2'C
	s := strings.TrimSpace(string(out))
	s = strings.TrimPrefix(s, "temp=")
	s = strings.TrimSuffix(s, "'C")

	val, _ := strconv.ParseFloat(s, 64)
	return val
}


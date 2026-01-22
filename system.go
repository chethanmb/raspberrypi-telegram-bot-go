package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func cpuTemp() string {
	out, err := exec.Command("vcgencmd", "measure_temp").Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out))
}

func uptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "N/A"
	}

	secStr := strings.Split(string(data), " ")[0]
	secs, _ := strconv.ParseFloat(secStr, 64)
	d := time.Duration(secs) * time.Second

	return fmt.Sprintf(
		"%dd %dh %dm",
		int(d.Hours())/24,
		int(d.Hours())%24,
		int(d.Minutes())%60,
	)
}

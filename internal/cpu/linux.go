//go:build linux
// +build linux

package cpu

import (
	"os/exec"
	"strconv"
	"strings"
)

func getInfo() CPUStats {
	cmd := exec.Command("top", "-b", "-n1")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	strOut := strings.Split(string(out), "\n")[2]
	cpuInfo := strings.Split(strings.Split(strOut, ":")[1], ",")
	resultStats := formCPUInfo(cpuInfo)

	return resultStats
}

func formCPUInfo(unformInfo []string) CPUStats {
	var us, sy, id, wa float32

	for _, item := range unformInfo {
		switch {
		case strings.Contains(item, "us"):
			us = paramToFloat(item)
		case strings.Contains(item, "sy"):
			sy = paramToFloat(item)
		case strings.Contains(item, "id"):
			id = paramToFloat(item)
		case strings.Contains(item, "wa"):
			wa = paramToFloat(item)
		}
	}
	return CPUStats{Sys: sy, Usr: us, Idle: id, Iowait: wa}
}

func paramToFloat(param string) float32 {
	param = strings.TrimSpace(param)
	strVal := strings.Split(param, " ")[0]
	floatVal, err := strconv.ParseFloat(strVal, 32)
	if err != nil {
		panic(err)
	}
	return float32(floatVal)
}

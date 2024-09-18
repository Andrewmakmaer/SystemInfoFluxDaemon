//go:build windows
// +build windows

package cpu

import (
	"encoding/binary"
	"log/slog"
	"math"
	"os/exec"
)

func getInfo() CPUStats {
	sys := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_Processor -Property PercentPrivilegedTime | Where-Object Name -eq '_Total').PercentPrivilegedTime`))
	usr := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_Processor -Property PercentUserTime | Where-Object Name -eq '_Total').PercentUserTime`))
	idle := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_Processor -Property PercentIdleTime | Where-Object Name -eq '_Total').PercentIdleTime`))
	iowait := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfDisk_PhysicalDisk | Where-Object Name -eq "_Total").PercentDiskTime`))

	return CPUStats{Sys: sys, Usr: usr, Idle: idle, Iowait: iowait}
}

func getFloatValue(cmd *exec.Cmd) float32 {
	val, err := cmd.Output()
	if err != nil {
		slog.Error("fail during search info", err)
	}
	bits := binary.LittleEndian.Uint32(val)
	float := math.Float32frombits(bits)
	return float
}

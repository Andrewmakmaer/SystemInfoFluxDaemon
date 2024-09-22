//go:build windows
// +build windows

package loadaverage

import (
	"encoding/binary"
	"log/slog"
	"math"
	"os/exec"
)

func getInfo() (result LoadInfo) {
	la1 := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_System).ProcessorQueueLength`))
	la5 := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_System).ProcessorQueueLength`))
	la15 := getFloatValue(exec.Command("powershell", "-Command",
		`(Get-CimInstance -ClassName Win32_PerfFormattedData_PerfOS_System).ProcessorQueueLength`))

	result = LoadInfo{La1: la1, La5: la5, La15: la15}
	return
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

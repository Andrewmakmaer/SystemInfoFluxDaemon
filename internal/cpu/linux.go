//go:build linux
// +build linux

package cpu

import (
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func getInfo() CPUStats {
	res, err := calculateValue()
	if err != nil {
		panic(err)
	}
	return res
}

func readCPUStats() ([]float32, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	cpuLine := strings.Fields(lines[0])[1:]

	var stats []float32
	for _, value := range cpuLine {
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		stats = append(stats, float32(v))
	}

	return stats, nil
}
func calculateValue() (CPUStats, error) {
	firstStats, err := readCPUStats()
	if err != nil {
		return CPUStats{}, err
	}

	time.Sleep(500 * time.Millisecond)

	secondStats, err := readCPUStats()
	if err != nil {
		return CPUStats{}, err
	}

	var firstTotal float32 = 0.0
	var secondTotal float32 = 0.0
	for i := range firstStats {
		firstTotal += firstStats[i]
		secondTotal += secondStats[i]
	}

	diff := secondTotal - firstTotal
	var result = CPUStats{
		Sys:    roundFloat((secondStats[2]-firstStats[2])/float32(diff)*100, 1),
		Usr:    roundFloat((secondStats[0]-firstStats[0])/float32(diff)*100, 1),
		Idle:   roundFloat((secondStats[3]-firstStats[3])/float32(diff)*100, 1),
		Iowait: roundFloat((secondStats[4]-firstStats[4])/float32(diff)*100, 1),
	}
	return result, nil
}

func roundFloat(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	res := float32(math.Round(float64(val)*ratio) / ratio)
	return res
}

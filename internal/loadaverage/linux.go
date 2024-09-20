//go:build linux
// +build linux

package loadaverage

import (
	"os/exec"
	"strconv"
	"strings"
)

func getInfo() (result LoadInfo) {
	cmd := exec.Command("cat", "/proc/loadavg")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	la := strings.Split(string(out), " ")
	result = LoadInfo{La1: paramToFloat(la[0]), La5: paramToFloat(la[1]), La15: paramToFloat(la[2])}

	return
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

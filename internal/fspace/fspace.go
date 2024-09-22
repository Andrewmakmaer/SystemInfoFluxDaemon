package fspace

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type FileSystemInfo struct {
	InodeUse int8
	SpaceUse int8
}

func GetInfo() map[string]FileSystemInfo {
	resultMap := make(map[string]FileSystemInfo)
	cmd := exec.Command("df", "--output=source,ipcent,pcent",
		"-x", "tmpfs",
		"-x", "overlay",
		"-x", "fuse.snapfuse",
		"-x", "9p")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	data := strings.Split(string(out), "\n")[1:]
	for _, item := range data {
		if item == "" {
			continue
		}
		r := regexp.MustCompile(`\s+`)
		replace := r.ReplaceAllString(item, " ")
		line := strings.Split(replace, " ")
		resultMap[line[0]] = FileSystemInfo{InodeUse: pcentToint(line[1]), SpaceUse: pcentToint(line[2])}
	}

	return resultMap
}

func pcentToint(pcent string) int8 {
	strRes := strings.Replace(pcent, "%", "", 1)
	intRes, err := strconv.Atoi(strRes)
	if err != nil {
		panic(err)
	}
	return int8(intRes)
}

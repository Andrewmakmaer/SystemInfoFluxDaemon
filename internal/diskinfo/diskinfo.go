package diskinfo

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type DiskInfo struct {
	Tps float32
	Kbs float32
}

func GetInfo() map[string]DiskInfo {
	resultMap := make(map[string]DiskInfo)
	cmd := exec.Command("iostat", "-d", "-k")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	data := strings.Split(string(out), "\n")[3:]
	for _, item := range data {
		if item == "" {
			continue
		}
		r := regexp.MustCompile(`\s+`)
		replace := r.ReplaceAllString(item, " ")
		line := strings.Split(replace, " ")

		tps, err := strconv.ParseFloat(line[1], 32)
		if err != nil {
			panic(err)
		}

		kbread, err := strconv.ParseFloat(line[2], 32)
		if err != nil {
			panic(err)
		}

		kbwrite, err := strconv.ParseFloat(line[3], 32)
		if err != nil {
			panic(err)
		}

		info := DiskInfo{Tps: float32(tps), Kbs: float32(kbread) + float32(kbwrite)}
		resultMap[line[0]] = info
	}

	return resultMap
}

package diskinfo

import (
	"log/slog"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type DiskInfo struct {
	Tps float32
	Kbs float32
}

func CheckRequirements() error {
	cmd := exec.Command("iostat")
	_, err := cmd.Output()
	if err.Error() == "exec: \"iostat\": executable file not found in $PATH" {
		slog.Error("not install iostat on host mashine")
		return err
	} else if err != nil {
		return err
	}
	return nil
}

func GetInfo() map[string]DiskInfo {
	resultMap := make(map[string]DiskInfo)
	cmd := exec.Command("iostat", "-d", "-k")
	out, err := cmd.Output()
	if err != nil {
		slog.Error(err.Error())
		return resultMap
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
			slog.Error("don't parse value to float", "value", line[1])
			return resultMap
		}

		kbread, err := strconv.ParseFloat(line[2], 32)
		if err != nil {
			slog.Error("don't parse value to float", "value", line[1])
			return resultMap
		}

		kbwrite, err := strconv.ParseFloat(line[3], 32)
		if err != nil {
			slog.Error("don't parse value to float", "value", line[1])
			return resultMap
		}

		info := DiskInfo{Tps: float32(tps), Kbs: float32(kbread) + float32(kbwrite)}
		resultMap[line[0]] = info
	}

	return resultMap
}

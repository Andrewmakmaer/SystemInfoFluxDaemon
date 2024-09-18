package cpu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormCpuInfo(t *testing.T) {
	var results []CPUStats

	for range 3 {
		results = append(results, GetInfo())
		time.Sleep(3 * time.Second)
	}

	t.Run("test correct result", func(t *testing.T) {
		item := results[0]
		check := item.Idle + item.Sys + item.Usr + item.Iowait
		require.LessOrEqual(t, check, float32(100))
		require.Less(t, float32(0), check)
	})

	t.Run("test what results is diferent", func(t *testing.T) {
		require.False(t, checkAllEqualing(results[0].Idle, results[1].Idle, results[2].Idle))
	})
}

func checkAllEqualing(values ...float32) bool {
	for i, item := range values {
		if i+1 <= len(values) {
			if item != values[i+1] {
				return false
			} else {
				continue
			}
		}
	}
	return true
}

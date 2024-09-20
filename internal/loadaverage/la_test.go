package loadaverage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormCpuInfo(t *testing.T) {
	t.Run("test correct result", func(t *testing.T) {
		for range 3 {
			check := GetInfo()
			require.Less(t, float32(0), check.La1)
			time.Sleep(3 * time.Second)
		}
	})
}

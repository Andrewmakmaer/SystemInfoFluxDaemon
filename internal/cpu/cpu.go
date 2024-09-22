package cpu

type CPUStats struct {
	Sys    float32
	Usr    float32
	Idle   float32
	Iowait float32
}

func GetInfo() CPUStats {
	return getInfo()
}

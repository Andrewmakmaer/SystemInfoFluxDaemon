package loadaverage

type LoadInfo struct {
	La1  float32
	La5  float32
	La15 float32
}

func GetInfo() (result LoadInfo) {
	return getInfo()
}

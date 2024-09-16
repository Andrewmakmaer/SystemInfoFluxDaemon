package app

import (
	"fmt"

	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/cpu"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/diskinfo"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/loadaverage"
)

type listCPU interface {
	GetRecords(count int) (result []any)
}

type listLa interface {
	GetRecords(count int) (result []any)
}

type listDisk interface {
	GetRecords(count int) (result []any)
}

type App struct {
	CPUStore  listCPU
	LaStore   listLa
	DiskStore listDisk
}

func (a *App) CPUValueAverage(m int) (sys, usr, idle, iowait float32, err error) {
	l := a.CPUStore.GetRecords(m)

	for _, item := range l {
		switch titem := item.(type) {
		case cpu.CPUStats:
			sys += titem.Sys
			usr += titem.Usr
			idle += titem.Idle
			iowait += titem.Iowait
		default:
			err = fmt.Errorf("unexpected type")
			return
		}
	}
	sys /= float32(m)
	usr /= float32(m)
	idle /= float32(m)
	iowait /= float32(m)
	return
}

func (a *App) LAverage(m int) (la1, la5, la15 float32, err error) {
	l := a.LaStore.GetRecords(m)

	for _, item := range l {
		switch titem := item.(type) {
		case loadaverage.LoadInfo:
			la1 += titem.La1
			la5 += titem.La5
			la15 += titem.La15
		default:
			err = fmt.Errorf("unexpected type")
			return
		}
	}
	la1 /= float32(m)
	la5 /= float32(m)
	la15 /= float32(m)
	return
}

func (a *App) DiskAverage(m int) (resultMap map[string]diskinfo.DiskInfo, err error) {
	l := a.DiskStore.GetRecords(m)
	resultMap = make(map[string]diskinfo.DiskInfo)

	for _, item := range l {
		switch titem := item.(type) {
		case map[string]diskinfo.DiskInfo:
			for d := range titem {
				resultMap = checkNilValueMap(resultMap, d)
				structItem := resultMap[d]
				structItem.Tps += titem[d].Tps
				structItem.Kbs += titem[d].Kbs

				resultMap[d] = structItem
			}
		default:
			err = fmt.Errorf("unexpected type")
			return
		}
	}
	for d := range resultMap {
		structItem := resultMap[d]
		structItem.Tps = resultMap[d].Tps / float32(m)
		structItem.Kbs = resultMap[d].Kbs / float32(m)

		resultMap[d] = structItem
	}
	return
}

func checkNilValueMap(m map[string]diskinfo.DiskInfo, v string) map[string]diskinfo.DiskInfo {
	_, ok := m[v]
	if !ok {
		item := diskinfo.DiskInfo{Tps: 0, Kbs: 0}
		m[v] = item
		return m
	}
	return m
}

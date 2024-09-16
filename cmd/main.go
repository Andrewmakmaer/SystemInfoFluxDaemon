package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/app"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/cpu"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/diskinfo"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/loadaverage"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/scheduler"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/server"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/storage/list"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()
	config := NewConfig(configFile)

	var (
		cpuList      = list.NewNodeList()
		loadAvgList  = list.NewNodeList()
		diskInfoList = list.NewNodeList()

		doingList []scheduler.Scheduler
	)

	for _, t := range config.Modes {
		switch t {
		case "cpu":
			doingList = append(doingList, scheduler.NewSchedule(1*time.Second, func() {
				cpuList.AddRecord(cpu.GetInfo())
			}))
		case "la":
			doingList = append(doingList, scheduler.NewSchedule(1*time.Second, func() {
				loadAvgList.AddRecord(loadaverage.GetInfo())
			}))
		case "disk":
			doingList = append(doingList, scheduler.NewSchedule(1*time.Second, func() {
				diskInfoList.AddRecord(diskinfo.GetInfo())
			}))
		}
	}

	application := app.App{
		CPUStore:  &cpuList,
		LaStore:   &loadAvgList,
		DiskStore: &diskInfoList,
	}
	srv := server.NewServer(&application, config.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		_, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			log.Fatal("failed to stop grps server: " + err.Error())
		}
	}()

	for _, task := range doingList {
		task.Do(ctx)
	}

	err := srv.Start(ctx)
	if err != nil {
		panic(err)
	}
}

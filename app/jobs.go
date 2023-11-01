package app

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/talkincode/esmqtt/common/zaplog/log"
)

var cronParser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

func (a *Application) initJob() {
	loc, _ := time.LoadLocation(a.appConfig.Location)
	a.sched = cron.New(cron.WithLocation(loc), cron.WithParser(cronParser))

	var err error
	_, err = a.sched.AddFunc("@every 30s", func() {
		go a.SchedSystemMonitorTask()
		go a.SchedProcessMonitorTask()
	})

	_, err = a.sched.AddFunc("@daily", func() {

	})

	if err != nil {
		log.Errorf("init job error %s", err.Error())
	}

	a.sched.Start()
}

// SchedSystemMonitorTask system monitor
func (a *Application) SchedSystemMonitorTask() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()
	//
	// timestamp := time.Now().Unix()
	//
	// var cpuuse float64
	// _cpuuse, err := cpu.Percent(0, false)
	// if err == nil && len(_cpuuse) > 0 {
	// 	cpuuse = _cpuuse[0]
	// }
	//
	// _meminfo, err := mem.VirtualMemory()
	// var memuse uint64
	// if err == nil {
	// 	memuse = _meminfo.Used
	// }

}

// SchedProcessMonitorTask app process monitor
func (a *Application) SchedProcessMonitorTask() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	// p, err := process.NewProcess(int32(os.Getpid()))
	// if err != nil {
	// 	return
	// }
	//
	// cpuuse, err := p.CPUPercent()
	// if err != nil {
	// 	cpuuse = 0
	// }
	//
	// if err != nil {
	// 	log.Error("add timeseries data error:", err.Error())
	// }
	//
	// meminfo, err := p.MemoryInfo()
	// if err != nil {
	// 	return
	// }
	// memuse := meminfo.RSS / 1024 / 1024
	//
	// if err != nil {
	// 	log.Error("add timeseries data error:", err.Error())
	// }
}

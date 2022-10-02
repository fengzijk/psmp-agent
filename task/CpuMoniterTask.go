package task

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
	"psmp-agent/cpu"
)

func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func InitTask() {

	log.Println("[Cron] Starting...")

	c := newWithSeconds()

	spec := viper.GetString("task-monitor-cron.cpu")

	// cpu监控
	_, _ = c.AddFunc(spec, func() {
		cpu.Monitor()
		log.Println("[Cron] Run cpuMonitor...")

	})

	//

	c.Start()
}

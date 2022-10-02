package task

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
	"psmp-agent/cpu"
	"psmp-agent/heartbeat"
)

func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func InitTask(ip string) {

	log.Println("[Cron] " + ip + "Starting...")

	c := newWithSeconds()

	cpuSpec := viper.GetString("task-monitor-cron.cpu")

	// cpu监控
	_, _ = c.AddFunc(cpuSpec, func() {
		cpu.Monitor()
		//log.Println("[Cron] Run cpuMonitor...")

	})

	// Agent 心跳
	heartbeatSpec := viper.GetString("task-monitor-cron.heartbeat")
	_, _ = c.AddFunc(heartbeatSpec, func() {
		heartbeat.AgentHeartbeat(ip)
		log.Println("[Cron] Run AgentHeartbeat...")
	})

	c.Start()
}

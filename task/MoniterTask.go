package task

import (
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"log"
	"psmp-agent/cpu"
	"psmp-agent/heartbeat"
	ipUtil "psmp-agent/ip"
)

func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func InitTask(ip string) {

	log.Println("[定时任务开启]" + ip + "  开始运行")

	c := newWithSeconds()

	// cpu监控
	if viper.GetBool("task-monitor-flag.cpuFlag") {
		cpuSpec := viper.GetString("task-monitor-cron.cpu")
		_, _ = c.AddFunc(cpuSpec, func() {
			cpu.Monitor()
			log.Println("cpu监控执行")

		})

	}
	// Agent 心跳
	if viper.GetBool("task-monitor-flag.diskFlag") {
		heartbeatSpec := viper.GetString("task-monitor-cron.heartbeatFlag")
		_, _ = c.AddFunc(heartbeatSpec, func() {
			heartbeat.AgentHeartbeat(ip)
			log.Println("心跳监控执行")
		})
	}
	// Agent 心跳
	if viper.GetBool("task-monitor-flag.heartbeatFlag") {
		ipChangeSpec := viper.GetString("task-monitor-cron.ipChange")
		_, _ = c.AddFunc(ipChangeSpec, func() {
			ipUtil.SendIPChange()
			log.Println("ip监控执行")
		})
	}

	c.Start()
}

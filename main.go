package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"psmp-agent/config"
	"psmp-agent/heartbeat"
	"psmp-agent/ip"
	"psmp-agent/task"
	"psmp-agent/util"
)

func main() {

	// 初始化配置
	initConfig()

	fmt.Println(util.GetDiskPercent())

	select {}
}

// initConfig
func initConfig() {

	//第一步 设置配置文件目录
	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("未找到application文件")
		} else {
			log.Println("读取application文件错误")
		}

		log.Println(err)
	}

	// cpu配置

	config.InitConfig()
	// 心跳配置
	heartbeat.InitConf()

	externalIP, _ := ip.ExternalIP()
	fmt.Println(externalIP)
	// 初始化定时任务
	task.InitTask(externalIP.String())
}

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"psmp-agent/ip"
	"psmp-agent/task"
)

func main() {

	// 初始化配置
	initConfig()
	// 初始化定时任务
	task.InitTask()

	//
	externalIP, _ := ip.ExternalIP()
	fmt.Println(externalIP)

	select {}
}

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
}
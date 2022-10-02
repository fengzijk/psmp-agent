package cpu

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"psmp-agent/ip"
	"psmp-agent/util"
)

const (
	// MonitorServerCpuSample  CPU监控样本前缀
	MonitorServerCpuSample = "server:cpu:sample:"
	// MonitorServerCpuAlarmInterval 告警间隔缓存key
	MonitorServerCpuAlarmInterval = "server:cpu:alarm_interval:"
	// MonitorServerCpuAlarmPre 是否上一次监控有告警缓存key
	MonitorServerCpuAlarmPre = "server:cpu:alarm_pre:"
)

func Monitor() {

	//样本list最大取样数量
	listMax := 30

	var cacheCpuList []float64

	externalIP, _ := ip.ExternalIP()

	// 样本缓存key
	sampleKey := MonitorServerCpuSample + externalIP.String()
	// 告警间隔缓存key
	alarmKey := MonitorServerCpuAlarmInterval + externalIP.String()
	// 是否上一次监控有告警缓存key
	preAlarmKey := MonitorServerCpuAlarmPre + externalIP.String()

	percent := util.GetCpuPercent()

	cache, b := util.GetCache(sampleKey)
	if b {
		err := json.Unmarshal([]byte(cache), &cacheCpuList)
		if err != nil {
			log.Print(err)
		}
	}
	cacheCpuList = append(cacheCpuList, percent)

	// 清理旧数据，只保留最新30条

	if len(cacheCpuList) > listMax {
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - 29):]
	}

	//cpu过载告警
	cpuOverloadAlarm(cacheCpuList, listMax, externalIP.String(), alarmKey, preAlarmKey, sampleKey)

	// cpu恢复正常通知
	cpuRecoveryNotification(cacheCpuList, externalIP.String(), preAlarmKey, sampleKey)

}

func cpuOverloadAlarm(cacheCpuList []float64, listMax int, ip, alarmKey, preAlarmKey, sampleKey string) {

	_, alarmBool := util.GetCache(alarmKey)

	//
	if len(cacheCpuList) > listMax && !alarmBool {
		// 清理旧数据，只保留最新30条
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - 29):]

		// 计算过载次数，即CPU使用率高于75%的次数
		overloadCount := 0
		for _, s := range cacheCpuList {
			if s > 75 {
				overloadCount++
			}
		}

		if float32(overloadCount) >= (float32(listMax) * 0.8) {
			// 发送告警邮件

			//

			// 告警缓存，frequency秒后失效，在此期间不会重复告警
			alarmCache := util.CacheModel{Key: alarmKey, Value: "1", ExpireSeconds: viper.GetInt("cpu.frequencySeconds")}
			util.SetCache(alarmCache)

			// 已告警标记，给恢复正常通知使用
			preAlarmCache := util.CacheModel{Key: preAlarmKey, Value: "1", ExpireSeconds: 100000}
			util.SetCache(preAlarmCache)
			doctorJson, _ := json.Marshal(cacheCpuList)
			// 存入本地List缓存
			sampleCache := util.CacheModel{Key: sampleKey, Value: string(doctorJson), ExpireSeconds: 100000}
			util.SetCache(sampleCache)
			return
		}

	}

}
func cpuRecoveryNotification(cacheCpuList []float64, ip, preAlarmKey, sampleKey string) {

	pr, preAlarmBool := util.GetCache(preAlarmKey)

	// 如果上一次有告警，且已恢复，则发恢复通知
	if len(cacheCpuList) >= 18 && preAlarmBool && "1" == pr {
		// 清理旧数据，只保留最新30
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - 29):]
		// 取最新18条（3分钟内）样本计算
		var ls []float64
		ls = cacheCpuList[(len(cacheCpuList) - 18):]

		// 计算过载次数，即CPU使用率高于75%的次数
		overloadCount := 0
		for _, s := range ls {
			if s > 75 {
				overloadCount++
			}
		}

		if overloadCount < 5 {
			// 发送邮件通知

			log.Print("发送邮件------------------------")
			// 已恢复告警标记
			preAlarmCache := util.CacheModel{Key: preAlarmKey, Value: "0", ExpireSeconds: 100000}
			util.SetCache(preAlarmCache)

		}
	}

	// 存入本地List缓存
	doctorJson, _ := json.Marshal(cacheCpuList)
	sampleCache := util.CacheModel{Key: sampleKey, Value: string(doctorJson), ExpireSeconds: 100000}
	util.SetCache(sampleCache)

}

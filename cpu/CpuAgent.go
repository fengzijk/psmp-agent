package cpu

import (
	"encoding/json"
	"log"
	"psmp-agent/config"
	"psmp-agent/util"
)

func Monitor() {

	//样本list最大取样数量
	listMax := config.CpuSampleSize

	var cacheCpuList []float64

	externalIP, _ := util.ExternalIP()

	// 样本缓存key
	sampleKey := config.MonitorServerCpuSample + externalIP.String()
	// 告警间隔缓存key
	alarmKey := config.MonitorServerCpuAlarmInterval + externalIP.String()
	// 是否上一次监控有告警缓存key
	preAlarmKey := config.MonitorServerCpuAlarmPre + externalIP.String()

	percent := util.GetCpuPercent()

	cache, b := util.GetCache(sampleKey)
	if b {
		err := json.Unmarshal([]byte(cache), &cacheCpuList)
		if err != nil {
			log.Print(err)
		}
	}
	cacheCpuList = append(cacheCpuList, percent)
	//
	//// 清理旧数据，只保留最新30条
	//
	//if len(cacheCpuList) > listMax {
	//	cacheCpuList = cacheCpuList[(len(cacheCpuList) - 29):]
	//}

	//cpu过载告警
	cpuOverloadAlarm(cacheCpuList, listMax, externalIP.String(), alarmKey, preAlarmKey, sampleKey)

	// cpu恢复正常通知
	cpuRecoveryNotification(cacheCpuList, externalIP.String(), preAlarmKey, sampleKey, listMax)

}

func cpuOverloadAlarm(cacheCpuList []float64, listMax int, ip, alarmKey, preAlarmKey, sampleKey string) {

	_, alarmBool := util.GetCache(alarmKey)

	//
	if len(cacheCpuList) > listMax && !alarmBool {
		// 清理旧数据，只保留最新30条
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - (listMax - 1)):]

		// 计算过载次数，即CPU使用率高于xx%的次数
		overloadCount := 0
		overloadThreshold := float64(config.CpuOverloadThreshold)
		for _, s := range cacheCpuList {
			if s > overloadThreshold {
				overloadCount++
			}
		}

		if float32(overloadCount) >= (float32(listMax) * 0.8) {

			// 告警缓存，frequency秒后失效，在此期间不会重复告警
			alarmCache := util.CacheModel{Key: alarmKey, Value: "1", ExpireSeconds: config.CpuFrequencySeconds}
			util.SetCache(alarmCache)

			// 已告警标记，给恢复正常通知使用
			preAlarmCache := util.CacheModel{Key: preAlarmKey, Value: "1", ExpireSeconds: 100000}
			util.SetCache(preAlarmCache)
			doctorJson, _ := json.Marshal(cacheCpuList)
			// 存入本地List缓存
			sampleCache := util.CacheModel{Key: sampleKey, Value: string(doctorJson), ExpireSeconds: 100000}
			util.SetCache(sampleCache)
			// 发送告警邮件
			util.NotifyEmailWebhook("psmp-agent", "", "", "CPU过载告警", ip+":CPU过载大于80%")
			return
		}

	}

}
func cpuRecoveryNotification(cacheCpuList []float64, ip, preAlarmKey, sampleKey string, listMax int) {

	pr, preAlarmBool := util.GetCache(preAlarmKey)

	// 如果上一次有告警，且已恢复，则发恢复通知
	if len(cacheCpuList) >= config.CpuNormalSampleSize && preAlarmBool && "1" == pr {
		// 清理旧数据，只保留最新30
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - (listMax - 1)):]
		// 取最新18条（3分钟内）样本计算
		var ls []float64
		ls = cacheCpuList[(len(cacheCpuList) - config.CpuNormalSampleSize):]

		// 计算过载次数，即CPU使用率高于75%的次数
		overloadCount := 0
		overloadThreshold := float64(config.CpuOverloadThreshold)
		for _, s := range ls {
			if s > overloadThreshold {
				overloadCount++
			}
		}

		if overloadCount < 5 {
			// 发送邮件通知
			util.NotifyEmailWebhook("psmp-agent", "", "", "CPU恢复正常", ip+":CPU恢复正常")

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

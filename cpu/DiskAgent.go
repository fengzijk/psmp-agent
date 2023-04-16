package cpu

import (
	"encoding/json"
	"log"
	"psmp-agent/config"
	"psmp-agent/ip"
	"psmp-agent/util"
)

func DiskMonitor() {

	//样本list最大取样数量
	listMax := config.DiskSampleSize

	var cacheCpuList []float64

	externalIP, _ := ip.ExternalIP()

	// 样本缓存key
	sampleKey := config.MonitorServerDiskSample + externalIP.String()
	// 告警间隔缓存key
	alarmKey := config.MonitorServerDiskAlarmInterval + externalIP.String()
	// 是否上一次监控有告警缓存key
	preAlarmKey := config.MonitorServerDiskAlarmPre + externalIP.String()

	percent := util.GetDiskPercent()

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

	//disk过载告警
	diskOverloadAlarm(cacheCpuList, listMax, externalIP.String(), alarmKey, preAlarmKey, sampleKey)

	// disk恢复正常通知
	diskRecoveryNotification(cacheCpuList, externalIP.String(), preAlarmKey, sampleKey, listMax)

}

func diskOverloadAlarm(cacheCpuList []float64, listMax int, ip, alarmKey, preAlarmKey, sampleKey string) {

	_, alarmBool := util.GetCache(alarmKey)

	//
	if len(cacheCpuList) > listMax && !alarmBool {
		// 清理旧数据，只保留最新30条
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - (listMax - 1)):]

		// 计算过载次数，即CPU使用率高于xx%的次数
		overloadCount := 0
		overloadThreshold := float64(config.DiskOverloadThreshold)
		for _, s := range cacheCpuList {
			if s > overloadThreshold {
				overloadCount++
			}
		}

		if float32(overloadCount) >= (float32(listMax) * 0.8) {
			// 发送告警邮件
			log.Print("发送邮件------------------------")
			var email = config.SendEmailRequest{FromName: "psmp-agent", Subject: "CPU过载告警", Body: ip + ":disk 过载 大于80%"}
			emailJson, _ := json.Marshal(email)

			_, err := util.PostJson(getEmailApi("disk"), string(emailJson), "")
			if err != nil {
				return
			}

			// 告警缓存，frequency秒后失效，在此期间不会重复告警
			alarmCache := util.CacheModel{Key: alarmKey, Value: "1", ExpireSeconds: config.DiskFrequencySeconds}
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
func diskRecoveryNotification(cacheCpuList []float64, ip, preAlarmKey, sampleKey string, listMax int) {

	pr, preAlarmBool := util.GetCache(preAlarmKey)

	// 如果上一次有告警，且已恢复，则发恢复通知
	if len(cacheCpuList) >= config.DiskNormalSampleSize && preAlarmBool && "1" == pr {
		// 清理旧数据，只保留最新30
		cacheCpuList = cacheCpuList[(len(cacheCpuList) - (listMax - 1)):]
		// 取最新18条（3分钟内）样本计算
		var ls []float64
		ls = cacheCpuList[(len(cacheCpuList) - config.DiskNormalSampleSize):]

		// 计算过载次数，即CPU使用率高于75%的次数
		overloadCount := 0
		overloadThreshold := float64(config.DiskOverloadThreshold)
		for _, s := range ls {
			if s > overloadThreshold {
				overloadCount++
			}
		}

		if overloadCount < 5 {
			// 发送邮件通知

			log.Print("发送邮件------------------------")

			var email = config.SendEmailRequest{FromName: "psmp-agent", Subject: "CPU过载告警", Body: ip + ":磁盘恢复正常"}
			emailJson, _ := json.Marshal(email)

			_, err := util.PostJson(getEmailApi("disk"), string(emailJson), "")
			if err != nil {
				return
			}

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

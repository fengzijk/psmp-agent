package config

import "github.com/spf13/viper"

const (
	// MonitorServerCpuSample  CPU监控样本前缀
	MonitorServerCpuSample = "server:cpu:sample:"
	// MonitorServerCpuAlarmInterval 告警间隔缓存key
	MonitorServerCpuAlarmInterval = "server:cpu:alarm_interval:"
	// MonitorServerCpuAlarmPre 是否上一次监控有告警缓存key
	MonitorServerCpuAlarmPre = "server:cpu:alarm_pre:"

	// MonitorServerDiskSample MonitorServerDisKSample 磁盘监控样本前缀
	MonitorServerDiskSample = "server:disk:sample:"
	// MonitorServerDiskAlarmInterval MonitorServerDisKAlarmInterval 磁盘告警间隔缓存key
	MonitorServerDiskAlarmInterval = "server:disk:alarm_interval:"
	// 磁盘 是否上一次监控有告警缓存key
	MonitorServerDiskAlarmPre = "server:disk:alarm_pre:"
)

type SendEmailRequest struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	FromName string `form:"fromName" json:"fromName"`
	ToUser   string `form:"ToUser" json:"toUser" `
	CcUser   string `form:"emailTo" json:"ccUser" `
	Subject  string `form:"subject" json:"subject" `
	Body     string `form:"content" json:"content" `
}

var (
	EmailUrl string

	CpuSampleSize int

	CpuOverloadThreshold int

	CpuFrequencySeconds int

	CpuNormalSampleSize int

	DiskSampleSize int

	DiskOverloadThreshold int

	DiskFrequencySeconds int

	DiskNormalSampleSize int
)

func InitConfig() {

	CpuSampleSize = viper.GetInt("cpu.sampleSize")
	CpuOverloadThreshold = viper.GetInt("cpu.overloadThreshold")
	CpuFrequencySeconds = viper.GetInt("frequencySeconds")
	CpuNormalSampleSize = viper.GetInt("cpu.normalSampleSize")

	DiskSampleSize = viper.GetInt("disk.sampleSize")
	DiskOverloadThreshold = viper.GetInt("disk.overloadThreshold")
	DiskFrequencySeconds = viper.GetInt("frequencySeconds")
	DiskNormalSampleSize = viper.GetInt("disk.normalSampleSize")

	EmailUrl = viper.GetString("psmp.url") + viper.GetString("psmp.send-email-api")

}

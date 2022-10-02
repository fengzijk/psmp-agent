package heartbeat

import (
	"fmt"
	"github.com/spf13/viper"
	"psmp-agent/util"
)

var Url string

func AgentHeartbeat(ip string) {

	url := fmt.Sprintf(Url, ip, "psmpAgent")
	util.GetJson(url, "")

}

func InitConf() {
	x := viper.GetString("psmp.url") + viper.GetString("psmp.heartbeat-api")
	Url = x
}

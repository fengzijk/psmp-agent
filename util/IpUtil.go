package util

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

const ipCacheShortKey = "ip_cache"

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetIpAddr2() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// 192.168.1.20:61085
	ip := strings.Split(localAddr.String(), ":")[0]

	return ip
}

func SendIPAlarm() {

	remoteUrlIP := GetRemoteUrlIP()

	ddnsIP, _ := GetDdnsIP()

	var atMobiles []string
	for _, tmp := range strings.Split(viper.GetString("ding-talk-webhook.atMobiles"), ";") {
		atMobiles = append(atMobiles, strings.TrimSpace(tmp))
	}

	var message = "【家庭路由器的公网IP监控】\n 出口公网IP:【" + remoteUrlIP + "】\n入口公网IP:【" + ddnsIP + "】"
	if remoteUrlIP != ddnsIP {
		message = "【家庭路由器的公网IP获取失败】-请检查是否公网\n" + "路由器IP:【" + ddnsIP + "】\n 外网IP:【" + remoteUrlIP + "】"
	}

	NotifyDingTalkWebhook(message, atMobiles)
	NotifyWeixinWebhook(message, atMobiles)
	NotifyEmailWebhook("【群晖】", "", "", "家庭路由器的公网IP监控】", message)
}

func GetPublicIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx]
}

func GetRemoteUrlIP() string {
	responseClient, errClient := http.Get("https://ipv4.netarm.com/") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		log.Print(errClient.Error())
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(responseClient.Body)

	body, _ := io.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))

	return clientIP
}

func GetDdnsIP() (string, error) {
	//addr, _ := net.ResolveIPAddr("ip", "fengzijk.synology.me")
	//conn, _ := net.Dial("ip:icmp", "fengzijk.synology.me")
	//
	//ip := conn.LocalAddr().String()
	//conn.Close()
	//a := strings.Split(ip, ":")
	//ip = a[0]
	//fmt.Println(ip)
	//fmt.Println(fmt.Sprintf("域名：%s 对应IP：%s 检测结果：正常 ", "fengzijk.synology.me", conn.LocalAddr().String()))

	var msg dns.Msg                               //创建一个Msg
	fqdn := dns.Fqdn(viper.GetString("ddns.url")) //调用fqdn将域转换为可以与DNS服务交换的FQDN
	msg.SetQuestion(fqdn, dns.TypeA)              //设置查询A记录
	in, err := dns.Exchange(&msg, "8.8.8.8:53")   //将消息发送到DNS服务器
	if err != nil {                               //判断是否有错误;如果有则打印输出
		return "", errors.New("获取域名dns地址失败")
	}
	if len(in.Answer) < 1 { //判断是否有响应内容,如果没有则输出没有记录并退出
		return "", errors.New("获取域名dns地址失败")
	}

	for _, answer := range in.Answer { //遍历所有应答
		if a, ok := answer.(*dns.A); ok { //将类型为A记录的类型取出;ok用于断言判断类型是否为*dns.A
			fmt.Println(a.A)               //
			return a.A.To4().String(), nil // 打印输出
		}
	}

	return "", errors.New("获取域名dns地址失败")

}

func SendIPChange() {

	var atMobiles []string
	for _, tmp := range strings.Split(viper.GetString("ding-talk-webhook.atMobiles"), ";") {
		atMobiles = append(atMobiles, strings.TrimSpace(tmp))
	}
	cacheKey := ipCacheShortKey

	// 缓存中查
	lastIp, _ := GetCache(cacheKey)

	remoteUrlIP := GetRemoteUrlIP()

	if len(lastIp) < 1 && len(remoteUrlIP) > 1 {
		ipCache := CacheModel{Key: cacheKey, Value: remoteUrlIP, ExpireSeconds: 60 * 60 * 4}
		SetCache(ipCache)
		return
	}

	ipCache := CacheModel{Key: cacheKey, Value: remoteUrlIP, ExpireSeconds: 60 * 60 * 4}
	SetCache(ipCache)
	var msg = "【家庭路由器的公网IP发生变化】\n由IP:【" + lastIp + "】\n变化为IP:【" + remoteUrlIP + "】"
	NotifyDingTalkWebhook(msg, atMobiles)
	NotifyWeixinWebhook(msg, atMobiles)
	NotifyEmailWebhook("【群晖】", "", "", "【家庭路由器的公网IP发生变化】", msg)
}

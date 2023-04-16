package util

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

type DingTalkWebhook struct {
	ActionCard struct {
		BtnOrientation string      `json:"btnOrientation"`
		Btns           interface{} `json:"btns"`
		HideAvatar     string      `json:"hideAvatar"`
		SingleTitle    string      `json:"singleTitle"`
		SingleURL      string      `json:"singleURL"`
		Text           string      `json:"text"`
		Title          string      `json:"title"`
	} `json:"actionCard"`
	At struct {
		AtMobiles interface{} `json:"atMobiles"`
		IsAtAll   bool        `json:"isAtAll"`
	} `json:"at"`
	FeedCard struct {
		Links interface{} `json:"links"`
	} `json:"feedCard"`
	Link struct {
		MessageURL string `json:"messageUrl"`
		PicURL     string `json:"picUrl"`
		Text       string `json:"text"`
		Title      string `json:"title"`
	} `json:"link"`
	Markdown struct {
		Text  string `json:"text"`
		Title string `json:"title"`
	} `json:"markdown"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type WeixinWebhook struct {
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content             string      `json:"content"`
		MentionedList       interface{} `json:"mentioned_list"`
		MentionedMobileList interface{} `json:"mentioned_mobile_list"`
	} `json:"text"`
}

type EmailWebhook struct {
	// binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
	FromName string `form:"fromName" json:"fromName"`
	ToUser   string `form:"ToUser" json:"toUser" `
	CcUser   string `form:"emailTo" json:"ccUser" `
	Subject  string `form:"subject" json:"subject" `
	Content  string `form:"content" json:"content" `
}

func NotifyDingTalkWebhook(msg string, atMobiles []string) {

	po := DingTalkWebhook{
		Text: struct {
			//Subject string `json:"subject"`
			Content string `json:"content"`
		}{Content: msg},
		At: struct {
			AtMobiles interface{} `json:"atMobiles"`
			IsAtAll   bool        `json:"isAtAll"`
		}{AtMobiles: atMobiles},
	}

	po.Msgtype = "text"

	messageJson, _ := json.Marshal(po)
	dingTalkUrl := fmt.Sprintf("%s?password=%s&dingSign=%s&dingToken=%s", viper.GetString("ding-talk-webhook-webhook.url"), viper.GetString("ding-talk-webhook.password"),
		viper.GetString("ding-talk-webhook.dingSign"), viper.GetString("ding-talk-webhook.dingToken"))
	_, _ = PostJson(dingTalkUrl, string(messageJson), "")

}

func NotifyWeixinWebhook(msg string, atMobiles []string) {
	po := WeixinWebhook{
		Msgtype: "text",
		Text: struct {
			Content             string      `json:"content"`
			MentionedList       interface{} `json:"mentioned_list"`
			MentionedMobileList interface{} `json:"mentioned_mobile_list"`
		}{Content: msg, MentionedMobileList: atMobiles},
	}

	messageJson, _ := json.Marshal(po)
	dingTalkUrl := fmt.Sprintf("%s?password=%s&weixinToken=%s", viper.GetString("weixin-webhook.url"), viper.GetString("weixin-webhook.password"),
		viper.GetString("weixin-webhook.weixinToken"))
	_, _ = PostJson(dingTalkUrl, string(messageJson), "")

}

func NotifyEmailWebhook(fromName, toUser, ccUser, subject, content string) {

	if len(toUser) < 1 {
		toUser = viper.GetString("email-webhook.toUser")
	}

	po := EmailWebhook{
		FromName: fromName,
		ToUser:   toUser,
		CcUser:   ccUser,
		Subject:  subject,
		Content:  content,
	}

	messageJson, _ := json.Marshal(po)

	dingTalkUrl := fmt.Sprintf("%s?password=%s", viper.GetString("email-webhook.url"), viper.GetString("email-webhook.password"))
	_, _ = PostJson(dingTalkUrl, string(messageJson), "")

}

func NotifyHeartbeat(fromName, toUser, ccUser, subject, content string) {

	if len(toUser) < 1 {
		toUser = viper.GetString("email-webhook.toUser")
	}

	po := EmailWebhook{
		FromName: fromName,
		ToUser:   toUser,
		CcUser:   ccUser,
		Subject:  subject,
		Content:  content,
	}

	messageJson, _ := json.Marshal(po)

	dingTalkUrl := fmt.Sprintf("%s?password=%s", viper.GetString("email-webhook.url"), viper.GetString("email-webhook.password"))
	_, _ = PostJson(dingTalkUrl, string(messageJson), "")

}

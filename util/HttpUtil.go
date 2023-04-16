package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func PostJson(url string, bodyJson string, authorization string) (string, error) {

	startTime := time.Now().UnixMilli()

	format := time.Now().Format("2006-01-02 03:04:05")

	contentType := "application/json"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyJson)))

	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorization))
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	all, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("响应码非200")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()

	}(resp.Body)

	if err != nil {
		log.Println(err)
		return "", err
	}
	res := string(all)
	endTime := time.Now().UnixMilli()
	log.Print(fmt.Sprintf("\n----------HTTPCLIENT请求开始------------\n请求时间   : %s \n请求方式   : %s\n请求URL    :%s  \n状态码     : %d \n结果       : %s \n耗时       : %d ms\n----------HTTPCLIENT请求结束------------\n",
		format,
		http.MethodGet,
		url,
		resp.StatusCode,
		res,
		endTime-startTime))

	return res, err
}

func GetJson(url string, authorization string) string {

	contentType := "application/json"
	startTime := time.Now().UnixMilli()

	format := time.Now().Format("2006-01-02 03:04:05")

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
		return ""
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authorization))
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}

	all, err := io.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if err != nil {
		log.Println(err)

		return ""
	}
	res := string(all)
	endTime := time.Now().UnixMilli()
	log.Print(fmt.Sprintf("\n----------HTTPCLIENT请求开始------------\n请求时间   : %s \n请求方式   : %s\n请求URL    :%s  \n状态码     : %d \n结果       : %s \n耗时       : %d ms\n----------HTTPCLIENT请求结束------------\n",
		format,
		http.MethodGet,
		url,
		resp.StatusCode,
		res,
		endTime-startTime))

	return res
}

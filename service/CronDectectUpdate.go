package service

import (
	"bytes"
	"fmt"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/service/tgService"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func cronDetectUpdate() {
	_, _ = DetectECodeUpdate()
}

//检测e码通网站更新
func DetectECodeUpdate() (bool, error) {
	loggerPrefix := "[service.cronDetectUpdate]"
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(1+rand.Intn(5)) * time.Second)

	localHtml, err := ioutil.ReadFile("./assets/ecode/index.html")
	if err != nil {
		logger.Error.Println(loggerPrefix + "读取本地页面失败:" + err.Error())
		return false, err
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", "https://ecode.neu.edu.cn/ecode/index.html", nil)
	if err != nil {
		logger.Error.Println(loggerPrefix + "创建请求失败:" + err.Error())
		return false, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh-TW;q=0.9,zh;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://ecode.neu.edu.cn/ecode/")
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"102\", \"Google Chrome\";v=\"102\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	rand.Seed(time.Now().UnixMilli())
	req.Header.Set("User-Agent", fmt.Sprintf("Mozilla/5.0 AppleWebKit/%d.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/%d.36", rand.Intn(300), rand.Intn(500)))

	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error.Println(loggerPrefix + "发送请求失败:" + err.Error())
		return false, err
	}
	defer resp.Body.Close()

	remoteHtml, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(loggerPrefix + "读取远程页面失败:" + err.Error())
		return false, err
	}

	//去除换行符
	localHtml = bytes.Replace(localHtml, []byte("\n"), []byte(""), -1)
	remoteHtml = bytes.Replace(remoteHtml, []byte("\n"), []byte(""), -1)

	if string(localHtml) != string(remoteHtml) {
		logger.Info.Println(loggerPrefix + "检测到e码通网站更新")
		tgService.SendTGMsg2Admin("检测到e码通网站更新")
		logger.Info.Println(loggerPrefix + "本地：" + string(localHtml))
		logger.Info.Println(loggerPrefix + "远程：" + string(remoteHtml))
		return true, nil
	} else {
		logger.Info.Println(loggerPrefix + "未检测到e码通网站更新")
		return false, nil
	}
}

package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
)

type RespAPITips struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
	Msg    string `json:"msg"`
}

func APITipsHandler(c *gin.Context) {
	md, err := ioutil.ReadFile("./assets/tips.md")
	if err != nil {
		utils.ReturnMsgJson(c, -1, "读取tips.md失败")
		return
	}

	//替换内容
	//[//]: # (copy right)
	md = bytes.Replace(md, []byte("https://enterneu.icu"), []byte(configs.Get().General.BaseUrl), -1)
	if utils.GetHostname() != "enterneu.icu" {
		md = bytes.Replace(md, []byte(`[//]: # "copy right"`), []byte("------\n> EnterNEU是一款按照[GPL-3.0 license](https://github.com/rroy233/EnterNEU/blob/main/LICENSE)协议开源的软件，完全免费。 © 2022 rroy233"), -1)
	}
	if configs.Get().TGService.Enabled == true {
		md = bytes.Replace(md, []byte(`[//]: # "your own bot instance"`), []byte(fmt.Sprintf("https://t.me/%s (本站实例)", configs.Get().TGService.BotUserName)), 1)
	}

	res := new(RespAPITips)
	res.Status = 0
	res.Data = string(md)

	//c.Data(200, "text/plain; charset=UTF-8", md)
	c.JSON(200, res)
	return
}

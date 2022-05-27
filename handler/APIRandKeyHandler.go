package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
	"strings"
)

type RespAPIRandKey struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		Key string `json:"key"`
		MD5 string `json:"MD5"`
	} `json:"data"`
}

func APIRandKeyHandler(c *gin.Context) {
	if configs.Get().General.WebEnabled == false {
		if configs.Get().TGService.Enabled == true {
			utils.ReturnMsgJson(c, -1, "暂时仅支持从Telegram Bot创建")
		} else {
			utils.ReturnMsgJson(c, -1, "服务已暂停")
		}
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		logger.Info.Println(err)
		utils.ReturnMsgJson(c, -1, "生成失败")
		return
	}
	resp := new(RespAPIRandKey)
	resp.Status = 0
	resp.Data.Key = strings.Replace(id.String(), "-", "", -1)
	resp.Data.MD5 = utils.MD5Short(resp.Data.Key + configs.Get().General.Md5Salt)
	c.JSON(200, resp)
	return
}

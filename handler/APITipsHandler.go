package handler

import (
	"github.com/gin-gonic/gin"
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

	res := new(RespAPITips)
	res.Status = 0
	res.Data = string(md)

	//c.Data(200, "text/plain; charset=UTF-8", md)
	c.JSON(200, res)
	return
}

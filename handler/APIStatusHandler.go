package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
	"time"
)

type RespAPIStatus struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		ImgUploaded bool `json:"ImgUploaded"`
		Info        struct {
			CheckUrl  string `json:"checkUrl"`
			StatusUrl string `json:"statusUrl"`
			DeleteUrl string `json:"deleteUrl"`
			ExpTime   string `json:"expTime"`
		} `json:"info"`
	} `json:"data"`
}

func APIStatusHandler(c *gin.Context) {
	token := c.Param("token")
	key := c.Param("key")

	helper := databases.NewHelper(key)
	helper.SetToken(token)
	if helper.Validate() == false {
		utils.ReturnMsgJson(c, -1, "token无效")
		return
	}

	store, err := helper.Get()
	if err != nil {
		utils.ReturnMsgJson(c, -1, "数据获取失败")
		return
	}

	res := new(RespAPIStatus)
	res.Status = 0
	res.Data.ImgUploaded = store.Student.ImgUploaded

	urlPrefix := fmt.Sprintf("%s/%s/%s", configs.Get().General.BaseUrl, token, key)
	res.Data.Info.CheckUrl = urlPrefix
	res.Data.Info.StatusUrl = urlPrefix + "/status"
	res.Data.Info.DeleteUrl = urlPrefix + "/delete"
	res.Data.Info.ExpTime = time.Unix(store.ExpTime, 0).Format("2006-01-02 15:04:05")
	c.JSON(200, res)
	return
}

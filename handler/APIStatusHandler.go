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
		Switch struct {
			Enabled bool   `json:"enabled"`
			Url     string `json:"url"`
		} `json:"switch"`
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

	res.Data.Info.CheckUrl = fmt.Sprintf("%s/api/%s/%s", utils.GetAPIBaseUrl(c), token, key)
	res.Data.Info.StatusUrl = fmt.Sprintf("%s/#/status/%s/%s", utils.GetFrontEndBaseUrl(c), token, key)
	res.Data.Info.DeleteUrl = fmt.Sprintf("%s/#/delete/%s/%s", utils.GetFrontEndBaseUrl(c), token, key)
	res.Data.Info.ExpTime = time.Unix(store.ExpTime, 0).Format("2006-01-02 15:04:05")

	//切换站点相关功能
	if configs.Get().Proxy.Enabled == true {
		res.Data.Switch.Enabled = true
		if configs.Get().Proxy.Enabled == true && c.GetHeader("X-API-Via") == configs.Get().Proxy.HeaderKey {
			res.Data.Switch.Url = configs.Get().General.BaseUrl
		} else {
			res.Data.Switch.Url = configs.Get().Proxy.FrontendBaseUrl
		}
	}

	c.JSON(200, res)
	return
}

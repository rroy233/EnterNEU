package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
	"time"
)

const (
	expTimeTypeOneHour = iota
	expTimeTypeOneDay
	expTimeTypeThreeDays
	expTimeTypeOneWeek
)

type ApiStoreReq struct {
	Key           string `json:"key" binding:"required"`
	KeyMD5        string `json:"keyMD5" binding:"required"`
	Name          string `json:"name" binding:"required"`
	StuID         string `json:"stuID" binding:"required"`
	EntranceName  string `json:"entranceName" binding:"required"`
	CodeType      string `json:"codeType" binding:"required"`
	ActualVehicle string `json:"actualVehicle" binding:"required"`
	ExpTimeType   int    `json:"expTimeType"`
}

type ApiStoreResp struct {
	Status int `json:"status"`
	Data   struct {
		VueUrl string `json:"vueUrl"`
	} `json:"data"`
}

func APIStoreHandler(c *gin.Context) {
	if configs.Get().General.WebEnabled == false {
		utils.ReturnMsgJson(c, -1, "生成失败")
		return
	}
	form := new(ApiStoreReq)
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Info.Println("[APIStoreHandler]参数无效：" + err.Error())
		utils.ReturnMsgJson(c, -1, "参数无效")
		return
	}

	//判断key是否有效
	if form.KeyMD5 != utils.MD5Short(form.Key+configs.Get().General.Md5Salt) {
		utils.ReturnMsgJson(c, -1, "key无效")
		return
	}

	//判断expDuration
	expDuration := 50 * time.Minute
	switch form.ExpTimeType {
	case expTimeTypeOneHour:
		expDuration = time.Hour
	case expTimeTypeOneDay:
		expDuration = 24 * time.Hour
	case expTimeTypeThreeDays:
		expDuration = 3 * 24 * time.Hour
	case expTimeTypeOneWeek:
		expDuration = 7 * 24 * time.Hour
	default:
		utils.ReturnMsgJson(c, -1, "ExpTimeType无效")
		return
	}

	helper := databases.NewHelper(form.Key)
	token, err := helper.CreateECode(form.Name, form.StuID, form.EntranceName, form.Key, "", form.CodeType, form.ActualVehicle, expDuration)
	if err != nil {
		utils.ReturnMsgJson(c, -1, "创建失败")
		return
	}

	resp := new(ApiStoreResp)
	resp.Status = 0
	resp.Data.VueUrl = fmt.Sprintf("/status/%s/%s", token, form.Key)

	c.JSON(200, resp)
	return
}

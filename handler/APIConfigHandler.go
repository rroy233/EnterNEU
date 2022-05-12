package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
	"strconv"
)

type RespAPIConfig struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		CodeType      []CodeConfigItems `json:"codeType"`
		ActualVehicle []CodeConfigItems `json:"actualVehicle"`
		ExpTime       []CodeConfigItems `json:"expTime"`
	} `json:"data"`
}

type CodeConfigItems struct {
	ConfigName  string `json:"configName"`
	ConfigValue string `json:"configValue"`
}

func APIConfigHandler(c *gin.Context) {
	ecodeConfig, err := configs.GetECodeConst()
	if err != nil {
		logger.Error.Println("读取配置失败" + err.Error())
		utils.ReturnMsgJson(c, -1, "读取配置失败")
		return
	}

	resp := new(RespAPIConfig)
	resp.Status = 0

	resp.Data.CodeType = make([]CodeConfigItems, 0)
	for i, codeType := range ecodeConfig.CodeTypes {
		resp.Data.CodeType = append(resp.Data.CodeType, CodeConfigItems{
			strconv.FormatInt(int64(i), 10),
			codeType,
		})
	}

	resp.Data.ActualVehicle = []CodeConfigItems{
		{ConfigName: "0", ConfigValue: "入"},
		{ConfigName: "1", ConfigValue: "出"},
	}
	resp.Data.ExpTime = []CodeConfigItems{
		{ConfigName: "0", ConfigValue: "1小时"},
		{ConfigName: "1", ConfigValue: "24小时"},
		{ConfigName: "2", ConfigValue: "3天"},
		{ConfigName: "3", ConfigValue: "一周"},
	}

	c.JSON(200, resp)
	return
}

package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
)

type eCodeUserPriorityResp struct {
	Data []eCodeUserPriorityData `json:"data"`
}

type eCodeUserPriorityData struct {
	Type       interface{} `json:"type"`
	Attributes struct {
		UserCode string `json:"userCode"`
		Ruleset  struct {
			Data struct {
				Id     string `json:"id"`
				Color  string `json:"color"`
				Name   string `json:"name"`
				Result string `json:"result"`
			} `json:"data"`
		} `json:"ruleset"`
	} `json:"attributes"`
}

func ECodeAPIUserPriorityHandler(c *gin.Context) {
	data, _ := c.Cookie("data")
	key, _ := c.Cookie("key")
	store := new(databases.StoreStruct)
	err := json.Unmarshal([]byte(data), store)
	if err != nil {
		c.JSON(200, ecodeUserInfoResp{})
		return
	}

	ec, err := configs.GetECodeConst()
	if err != nil {
		logger.Error.Println("[ECodeAPIUserPriorityHandler]读取配置文件发生错误", err.Error())
		utils.ReturnMsgJson(c, -1, "读取配置文件发生错误")
		return
	}

	res := new(eCodeUserPriorityResp)
	res.Data = make([]eCodeUserPriorityData, 1)
	res.Data[0].Type = struct{}{}
	res.Data[0].Attributes.UserCode, _ = utils.AesCbcDecryptString(store.Student.ID, key)
	res.Data[0].Attributes.Ruleset.Data.Id = "114514"
	res.Data[0].Attributes.Ruleset.Data.Result = store.Entrance.CodeType
	res.Data[0].Attributes.Ruleset.Data.Color = ec.ColorByCodeTypeID[store.Entrance.CodeType]
	res.Data[0].Attributes.Ruleset.Data.Name = ec.CodeTypeTextByIndex[store.Entrance.CodeType]

	c.JSON(200, res)
	return
}

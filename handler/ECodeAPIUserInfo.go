package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
)

type ecodeUserInfoResp struct {
	Data []ecodeUserInfoData `json:"data"`
}

type ecodeUserInfoData struct {
	Type       interface{} `json:"type"`
	Attributes struct {
		UserCode string `json:"userCode"`
		UserName string `json:"userName"`
	} `json:"attributes"`
}

func ECodeAPIUserInfo(c *gin.Context) {
	data, _ := c.Cookie("data")
	key, _ := c.Cookie("key")
	store := new(databases.StoreStruct)
	err := json.Unmarshal([]byte(data), store)
	if err != nil {
		c.JSON(200, ecodeUserInfoResp{})
		return
	}

	info := new(ecodeUserInfoResp)
	info.Data = make([]ecodeUserInfoData, 1)
	info.Data[0].Type = struct{}{}
	info.Data[0].Attributes.UserCode, _ = utils.AesCbcDecryptString(store.Student.ID, key)
	info.Data[0].Attributes.UserName, _ = utils.AesCbcDecryptString(store.Student.Name, key)
	c.JSON(200, info)
}

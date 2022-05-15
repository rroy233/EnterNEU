package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
	"strconv"
	"time"
)

type ecodeRegisterResp struct {
	Data []ecodeRegisterData `json:"data"`
}

type ecodeRegisterData struct {
	Attributes struct {
		QrCreateTime string `json:"qrCreateTime"`
		Ruleset      struct {
			Color  string `json:"color"`
			Result string `json:"result"`
			Name   string `json:"name"`
			Id     string `json:"id"`
		} `json:"ruleset"`
		ActualVehicle string `json:"actualVehicle"`
		UserCode      string `json:"userCode"`
		Result        string `json:"result"`
		CreateTime    string `json:"createTime"`
		Id            string `json:"id"`
		Entrance      string `json:"entrance"`
	} `json:"attributes"`
}

func ECodeAPIRegister(c *gin.Context) {
	data, _ := c.Cookie("data")
	key, _ := c.Cookie("key")
	store := new(databases.StoreStruct)
	err := json.Unmarshal([]byte(data), store)
	if err != nil {
		c.JSON(200, ecodeRegisterResp{})
		return
	}

	entranceName := store.Entrance.Name
	entranceId := 444

	info := new(ecodeRegisterResp)
	info.Data = make([]ecodeRegisterData, 1)
	/*
		前端代码
		t.userCode = e.data[0].attributes.userCode,
		t.createTime = e.data[0].attributes.createTime,
		t.entrance = JSON.parse(e.data[0].attributes.entrance).data.name,
		t.actualVehicle = e.data[0].attributes.actualVehicle % 2 === 0 ? "入": "出",
		t.passText = Z(e.data[0].attributes.result),
		t.passColor = e.data[0].attributes.ruleset.color
	*/
	color := "#509674"
	codeType := store.Entrance.CodeType
	ecodeConfigs, err := configs.GetECodeConst()
	if err == nil {
		id, _ := strconv.Atoi(store.Entrance.CodeType)
		color = ecodeConfigs.Colors[id]
	}
	info.Data[0].Attributes.QrCreateTime = time.Now().Format("2006-01-02T15:04:05")
	info.Data[0].Attributes.Ruleset.Color = color
	info.Data[0].Attributes.Ruleset.Name = "同意入校"
	info.Data[0].Attributes.Ruleset.Id = "114514"
	info.Data[0].Attributes.Ruleset.Result = "1"
	info.Data[0].Attributes.ActualVehicle = store.Entrance.ActualVehicle //0:入，1:出
	info.Data[0].Attributes.UserCode, _ = utils.AesCbcDecryptString(store.Student.ID, key)
	info.Data[0].Attributes.Result = codeType //前端由此判断文字
	info.Data[0].Attributes.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	info.Data[0].Attributes.Id = "114514"

	info.Data[0].Attributes.Entrance = fmt.Sprintf(`{"data":{"name":"%s","id":%d}}`, entranceName, entranceId)

	c.JSON(200, info)
}

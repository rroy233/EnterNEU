package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/rroy233/EnterNEU/logger"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
)

type ecodeUserFaceResp struct {
	Data []ecodeUserFaceData `json:"data"`
}

type ecodeUserFaceData struct {
	Type       interface{} `json:"type"`
	Attributes struct {
		FaceImg  string `json:"faceImg"`
		UserCode string `json:"userCode"`
	} `json:"attributes"`
}

func ECodeAPIUserFace(c *gin.Context) {
	data, _ := c.Cookie("data")
	key, _ := c.Cookie("key")
	store := new(databases.StoreStruct)
	err := json.Unmarshal([]byte(data), store)
	if err != nil {
		c.JSON(200, ecodeRegisterResp{})
		return
	}
	info := new(ecodeUserFaceResp)
	info.Data = make([]ecodeUserFaceData, 1)
	info.Data[0].Type = struct{}{}
	if store.Student.ImgUploaded == false {
		info.Data[0].Attributes.FaceImg = ""
	} else {
		imgData, err := ioutil.ReadFile(store.Student.ImgPath)
		if err != nil {
			logger.Info.Println("读取失败：", err)
		} else {
			de, err := utils.AesCbcDecrypt(imgData, key)
			if err != nil {
				logger.Info.Println("解密失败：", err)
			} else {
				info.Data[0].Attributes.FaceImg = "data:" + store.Student.ImgContentType + ";base64," + base64.StdEncoding.EncodeToString(de)
			}

		}
	}
	info.Data[0].Attributes.UserCode = "114514"
	c.JSON(200, info)
}

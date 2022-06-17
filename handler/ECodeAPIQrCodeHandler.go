package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type ecodeQrcodeResp struct {
	Data []ecodeQrcodeData `json:"data"`
}

type ecodeQrcodeData struct {
	Type       interface{} `json:"type"`
	Attributes struct {
		QrCode        string `json:"qrCode"`
		CreateTime    string `json:"createTime"`
		QrInvalidTime string `json:"qrInvalidTime"`
	} `json:"attributes"`
}

func ECodeAPIQrCodeHandler(c *gin.Context) {
	res := new(ecodeQrcodeResp)
	res.Data = make([]ecodeQrcodeData, 1)
	res.Data[0].Type = struct{}{}
	res.Data[0].Attributes.CreateTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
	res.Data[0].Attributes.QrCode = "NEU" + strconv.FormatInt(time.Now().UnixNano(), 10) + strings.Repeat("0", 5)
	res.Data[0].Attributes.QrInvalidTime = strconv.FormatInt(time.Now().Add(5*time.Minute).UnixMilli(), 10)
	c.JSON(200, res)
	return
}

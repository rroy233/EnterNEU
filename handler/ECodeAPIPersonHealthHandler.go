package handler

import (
	"github.com/gin-gonic/gin"
)

type ecodePersionHealthResp struct {
	Data struct {
		Type       interface{} `json:"type"`
		Attributes struct {
			Records []ecodePersionHealthRecords `json:"records"`
		} `json:"attributes"`
	} `json:"data"`
}

type ecodePersionHealthRecords struct {
	CreatedOn   string `json:"createdOn"`
	CreatedAt   string `json:"createdAt"`
	HealthState string `json:"healthState"`
}

func ECodeAPIPersonHealthHandler(c *gin.Context) {
	res := new(ecodePersionHealthResp)
	res.Data.Type = struct{}{}
	res.Data.Attributes.Records = make([]ecodePersionHealthRecords, 1)
	res.Data.Attributes.Records[0].CreatedAt = "2022-01-14T05:14:22.000000Z"
	res.Data.Attributes.Records[0].HealthState = "正常"
	res.Data.Attributes.Records[0].CreatedOn = "2022-01-14"
	c.JSON(200, res)
	return
}

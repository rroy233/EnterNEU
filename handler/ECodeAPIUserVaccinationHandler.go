package handler

import "github.com/gin-gonic/gin"

func ECodeAPIUserVaccinationHandler(c *gin.Context) {
	c.Data(200, "application/json;charset=utf-8", []byte(`{"data":[{"type":null,"attributes":{"dose":"1","date":"2021-05-22"}},{"type":null,"attributes":{"dose":"2","date":"2021-06-22"}},{"type":null,"attributes":{"dose":"3","date":"2021-07-22"}}]}`))
	return
}

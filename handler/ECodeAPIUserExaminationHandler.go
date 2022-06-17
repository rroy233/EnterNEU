package handler

import "github.com/gin-gonic/gin"

func ECodeAPIUserExaminationHandler(c *gin.Context) {
	c.Data(200, "application/json;charset=utf-8", []byte(`{"data":{"type":null,"attributes":{"records":[{"id":"null","userCode":"114514","expiringAt":"null"}]}}}`))
	return
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
)

func ECodeIndexHandler(c *gin.Context) {
	data, err := ioutil.ReadFile("./assets/ecode/index.html")
	if err != nil {
		utils.ReturnPlainHtml(c, "找不到./assets/ecode/index.html")
		return
	}
	c.Data(200, gin.MIMEHTML, data)
	return
}

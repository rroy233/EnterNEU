package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
)

func IndexHandler(c *gin.Context) {
	index, err := ioutil.ReadFile("./assets/enterneu/index.html")
	if err != nil {
		utils.ReturnPlainHtml(c, "模板读取失败")
		return
	}

	c.Data(200, "text/html;charset=utf-8", index)
	return
}

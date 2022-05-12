package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
)

func APIDeleteHandler(c *gin.Context) {
	token := c.Param("token")
	key := c.Param("key")

	helper := databases.NewHelper(key)
	helper.SetToken(token)
	if helper.Validate() == false {
		utils.ReturnMsgJson(c, -1, "token或key无效")
		return
	}

	err := helper.Delete()
	if err != nil {
		utils.ReturnMsgJson(c, -1, "删除失败")
		return
	}

	utils.ReturnMsgJson(c, 0, "已删除")
	return
}

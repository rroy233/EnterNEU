package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
	"strings"
	"time"
)

func TokenShadowrocketHandler(c *gin.Context) {
	token := c.Param("token")
	key := c.Param("key")

	helper := databases.NewHelper(key)
	helper.SetToken(token)
	if helper.Validate() == false {
		utils.ReturnMsgJson(c, -1, "token无效")
		return
	}

	store, err := helper.Get()
	if err != nil {
		utils.ReturnMsgJson(c, -1, "获取数据失败")
		return
	}

	confFile, err := ioutil.ReadFile("./configs/shadowrocket.conf")
	if err != nil {
		utils.ReturnMsgJson(c, -1, "模板读取失败")
		return
	}

	fileName := fmt.Sprintf("%s_%s过期.conf", token, time.Unix(store.ExpTime, 0).Format("01月02日15:04"))
	Url := fmt.Sprintf("%s/%s/%s?_=", utils.GetAPIBaseUrl(c), token, key)

	//替换
	fileData := strings.Replace(string(confFile), "{{url_replace}}", Url, -1)

	c.Header("content-disposition", "attachment; filename="+fileName)
	c.Data(200, "application/octet-stream; charset=utf-8", []byte(fileData))
	return
}

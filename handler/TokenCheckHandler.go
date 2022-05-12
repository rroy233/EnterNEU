package handler

import (
	"encoding/json"
	"fmt"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/utils"
)

func TokenCheckHandler(c *gin.Context) {
	token := c.Param("token")
	key := c.Param("key")

	helper := databases.NewHelper(key)
	helper.SetToken(token)
	data, err := helper.Get()
	if err != nil {
		logger.Info.Println(err)
		utils.ReturnMsgJson(c, -1, "参数无效")
		return
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		logger.Info.Println(err)
		utils.ReturnMsgJson(c, -1, "系统异常")
		return
	}

	c.SetCookie("data", string(dataJson), 60*60, "/", utils.GetHostname(), utils.GetCookieSecure(), true)
	c.SetCookie("key", key, 60*60, "/", utils.GetHostname(), utils.GetCookieSecure(), true)
	c.SetCookie("exp_time", strconv.FormatInt(data.ExpTime, 10), 60*60, "/", utils.GetHostname(), utils.GetCookieSecure(), true)
	c.SetCookie("checksum", utils.MD5Short(string(dataJson)+strconv.FormatInt(data.ExpTime, 10)+key), 60*60, "/", utils.GetHostname(), utils.GetCookieSecure(), true)
	c.Redirect(302, fmt.Sprintf("%s/ecode/#/codeRecord?entrance_id=1&actual_vehicle=0", configs.Get().General.BaseUrl))
	return
}

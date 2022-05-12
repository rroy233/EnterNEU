package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/utils"
)

func ECodeMiddleWare(c *gin.Context) {
	data, err := c.Cookie("data")
	checksum, err := c.Cookie("checksum")
	expTime, err := c.Cookie("exp_time")
	key, err := c.Cookie("key")
	if err != nil {
		utils.ReturnPlainHtml(c, "请从生成的地址进入。")
		c.Abort()
		return
	}

	//校验和
	if utils.MD5Short(data+expTime+key) != checksum {
		utils.ReturnPlainHtml(c, "您当前的cookie无效")
		c.Abort()
		return
	}

	//判断是否过期
	exp, err := strconv.ParseInt(expTime, 10, 64)
	if err != nil {
		utils.ReturnPlainHtml(c, "您当前的cookie无效")
		c.Abort()
		return
	}

	if exp < time.Now().Unix() {
		utils.ReturnPlainHtml(c, "您的信息已过期，请重新生成")
		c.Abort()
		return
	}

	c.Next()
}

package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
)

func TokenStatusHandler(c *gin.Context) {
	c.Redirect(302, fmt.Sprintf("%s/#/status/%s/%s", configs.Get().General.BaseUrl, c.Param("token"), c.Param("key")))
	return
}

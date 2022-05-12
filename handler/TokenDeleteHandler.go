package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func TokenDeleteHandler(c *gin.Context) {
	c.Redirect(200, fmt.Sprintf("/#/delete/%s/%s", c.Param("token"), c.Param("key")))
	return
}

package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
)

func CorsMiddleware(c *gin.Context) {
	if configs.Get().General.Production == false {
		c.Header("Access-Control-Allow-Origin", "*")
	}
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "content-type")
		c.Status(200)
		return
	}
	c.Next()
	return
}

package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func CacheMiddleware(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.Next()
		return
	}
	ext := []string{"gif", "jpg", "png", "js", "css", "js.map"}
	has := false
	for _, s := range ext {
		if strings.HasSuffix(c.Request.URL.Path, s) == true {
			c.Header("cache-control", "max-age=43201")
			has = true
			break
		}
	}
	if has == false {
		c.Header("cache-control", "no-cache")
	}
	c.Next()
}

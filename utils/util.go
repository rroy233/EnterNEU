package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rroy233/EnterNEU/configs"
	"net/url"
	"os"
	"strings"
)

func ReturnPlainHtml(c *gin.Context, text string) {
	c.Data(200, "text/html; charset=UTF-8", []byte(text))
}

func ReturnMsgJson(c *gin.Context, status int, msg string) {
	c.JSON(200, &gin.H{
		"status": status,
		"msg":    msg,
		"data":   struct{}{},
	})
}

func Ctx() context.Context {
	return context.Background()
}

// MD5Short 生成6位MD5
func MD5Short(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil)[0:5])
}

func FsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func LogGetIP(param *gin.LogFormatterParams) string {
	if param.Request.Header.Get("X-Forwarded-For") != "" {
		return param.Request.Header.Get("X-Forwarded-For")
	} else {
		return param.ClientIP
	}
}

func LogGetPath(param *gin.LogFormatterParams) string {
	part := strings.Split(param.Path[1:], "/")
	path := ""
	if len(part) >= 3 && len(part[1]) == 10 && len(part[2]) == 32 { //包含token和key
		part[1] = part[1][:5] + strings.Repeat("*", 10-5)
		part[2] = part[2][:10] + strings.Repeat("*", 32-10)
		path = "/" + strings.Join(part, "/")
	} else {
		path = param.Path
	}
	return path
}

func NewUUIDToken() string {
	u, _ := uuid.NewUUID()
	return strings.Replace(u.String(), "-", "", -1)
}

func Sha256Hex(data []byte) string {
	digest := sha256.New()
	digest.Write(data)
	return hex.EncodeToString(digest.Sum(nil))
}

func GetHostname(c *gin.Context) string {
	base := ""
	if configs.Get().Proxy.Enabled == true && c.GetHeader("X-API-Via") == configs.Get().Proxy.HeaderKey {
		base = configs.Get().Proxy.ApiBaseUrl
	} else {
		base = configs.Get().General.BaseUrl
	}
	u, err := url.Parse(base)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

func GetCookieSecure() bool {
	u, err := url.Parse(configs.Get().General.BaseUrl)
	if err != nil {
		return false
	}
	if u.Scheme == "https" {
		return true
	} else {
		return false
	}
}

func GetAPIBaseUrl(c *gin.Context) string {
	if configs.Get().Proxy.Enabled == true && c.GetHeader("X-API-Via") == configs.Get().Proxy.HeaderKey {
		return configs.Get().Proxy.ApiBaseUrl
	}
	return configs.Get().General.BaseUrl
}

func GetFrontEndBaseUrl(c *gin.Context) string {
	if configs.Get().Proxy.Enabled == true && c.GetHeader("X-API-Via") == configs.Get().Proxy.HeaderKey {
		return configs.Get().Proxy.FrontendBaseUrl
	}
	return configs.Get().General.BaseUrl
}

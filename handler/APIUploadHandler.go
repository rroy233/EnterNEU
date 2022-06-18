package handler

import (
	"fmt"
	"github.com/rroy233/EnterNEU/logger"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/utils"
)

var allowContentType map[string]int = map[string]int{
	"image/jpeg": 1,
	"image/png":  1,
}

func APIUploadHandler(c *gin.Context) {
	token := c.Query("token")
	key := c.Query("key")

	if token == "" || key == "" {
		utils.ReturnMsgJson(c, -1, "参数无效")
		return
	}

	helper := databases.NewHelper(key)
	helper.SetToken(token)
	if helper.Validate() == false {
		utils.ReturnMsgJson(c, -1, "token无效")
		return
	}

	//读取文件
	file, err := c.FormFile("file")
	if err != nil {
		logger.Info.Println("[文件上传]文件装载失败", err)
		utils.ReturnMsgJson(c, -1, fmt.Sprintf("文件上传失败(%s)", err.Error()))
		return
	}

	//文件大小
	sizeLimit := int64(configs.MaxUploadSize)
	if file.Size > sizeLimit {
		utils.ReturnMsgJson(c, -1, "文件大小超过规定值")
		return
	}

	contentType := file.Header.Values("Content-Type")[0]

	//安全检查
	if allowContentType[contentType] != 1 {
		utils.ReturnMsgJson(c, -1, "文件格式无效")
		return
	}
	if strings.Contains(file.Filename, "..") {
		utils.ReturnMsgJson(c, -1, "文件名无效")
		return
	}

	//保存临时文件
	tmpFileName := utils.MD5Short(fmt.Sprintf("tmpfile_%d", time.Now().UnixMilli()))
	err = c.SaveUploadedFile(file, "./storage/tmp/"+tmpFileName)
	if err != nil {
		utils.ReturnMsgJson(c, -1, err.Error())
		return
	}
	defer func() {
		//清除临时文件
		_ = os.Remove("./storage/tmp/" + tmpFileName)
	}()

	//读取临时文件
	imgo, err := ioutil.ReadFile("./storage/tmp/" + tmpFileName)
	if err != nil {
		utils.ReturnMsgJson(c, -1, err.Error())
		return
	}

	//加密头像，存储
	imgEncrypt, err := utils.AesCbcEncrypt(imgo, key)
	err = ioutil.WriteFile(fmt.Sprintf("./storage/%s.data", token), imgEncrypt, 0755)
	if err != nil {
		utils.ReturnMsgJson(c, -1, err.Error())
		return
	}

	//更新redis
	err = helper.UpdateImg(fmt.Sprintf("./storage/%s.data", token), contentType)
	if err != nil {
		utils.ReturnMsgJson(c, -1, err.Error())
		return
	}

	c.JSON(200, gin.H{"status": 0})
	return
}

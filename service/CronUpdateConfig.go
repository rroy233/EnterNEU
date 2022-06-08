package service

import (
	"bytes"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
	"io/ioutil"
)

var oldConfig []byte
var configPath = "./config.yaml"

func cronUpdateConfig() {
	loggerPrefix := "[service.cronUpdateConfig]"
	fileData, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
		return
	}
	if oldConfig == nil {
		oldConfig = fileData
	} else {
		if bytes.Equal(fileData, oldConfig) == false {
			configs.InitConfig(configPath)
			logger.Info.Println(loggerPrefix + "已自动更新配置文件")
			oldConfig = fileData
			return
		}
	}
}

package configs

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	General struct {
		Production bool   `yaml:"production"`
		BaseUrl    string `yaml:"baseUrl"`
		ListenPort string `yaml:"listenPort"`
		AesIv      string `yaml:"aesIv"`
		Md5Salt    string `yaml:"md5Salt"`
	} `yaml:"general"`
	TGService struct {
		Enabled         bool   `yaml:"enabled"`
		BotToken        string `yaml:"botToken"`
		AdminUID        int64  `yaml:"adminUID"`
		HandleWorkerNum int    `yaml:"handleWorkerNum"`
		MsgSenderNum    int    `yaml:"msgSenderNum"`
	} `yaml:"tgService"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"DB"`
	} `yaml:"redis"`
}

var config *ConfigStruct

func InitConfig(path string) {
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("配置文件加载失败！", err)
	}
	err = yaml.Unmarshal(confFile, &config)
	if err != nil {
		log.Fatalln("配置文件加载失败！(" + err.Error() + ")")
	}
	return
}

func Get() *ConfigStruct {
	return config
}

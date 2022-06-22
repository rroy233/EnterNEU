package configs

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	General struct {
		Production       bool   `yaml:"production"`
		WebEnabled       bool   `yaml:"web_enabled"`
		BaseUrl          string `yaml:"baseUrl"`
		ListenPort       string `yaml:"listenPort"`
		AesIv            string `yaml:"aesIv"`
		Md5Salt          string `yaml:"md5Salt"`
		AutoDetectUpdate bool   `yaml:"autoDetectUpdate"`
		AutoDetectTime   string `yaml:"autoDetectTime"`
	} `yaml:"general"`
	TGService struct {
		Enabled         bool   `yaml:"enabled"`
		BotToken        string `yaml:"botToken"`
		BotUserName     string `yaml:"botUserName"`
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
	StatisticsReport struct {
		V651La    bool   `yaml:"v6_51_la"`
		V651LaJs1 string `yaml:"v6_51_la_js1"`
		V651LaJs2 string `yaml:"v6_51_la_js2"`
	} `yaml:"statistics_report"`
	Proxy struct {
		Enabled         bool   `yaml:"enabled"`
		FrontendBaseUrl string `yaml:"frontend_baseurl"`
		ApiBaseUrl      string `yaml:"api_baseurl"`
		HeaderKey       string `yaml:"header_key"`
	} `yaml:"proxy"`
}

var config *ConfigStruct

const MaxUploadSize = 1 << 20 //1MiB

func InitConfig(path string) {
	config = new(ConfigStruct)
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("配置文件加载失败！", err)
	}
	err = yaml.Unmarshal(confFile, &config)
	if err != nil {
		log.Fatalln("配置文件加载失败！(" + err.Error() + ")")
	}

	//检查项
	if config.General.AutoDetectUpdate == true && config.TGService.Enabled == false {
		log.Fatalln("自动检测更新需要开启[TG服务]，并填写好管理员uid。")
	}

	return
}

func Get() *ConfigStruct {
	return config
}

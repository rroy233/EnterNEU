package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/router"
	"github.com/rroy233/EnterNEU/service"
	"github.com/rroy233/EnterNEU/service/tgService"
	"github.com/rroy233/EnterNEU/utils"
	"io"
	"os"
)

func main() {
	r := gin.New()

	//读取配置
	configs.InitConfig("./config.yaml")

	//初始化日志
	logger.New(configs.Get().General.Production)

	//gin日志输出
	logFile, err := os.OpenFile("./log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.FATAL.Fatalln(err)
	}
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s]- %s - \"%s %s %d %s \"%s\" %s\"\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				utils.LogGetIP(&param),
				param.Method,
				//param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		Output:    io.MultiWriter(logFile),
		SkipPaths: []string{},
	}))
	r.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)

	//文件上传大小限制
	r.MaxMultipartMemory = 5 << 20 // 5 MB

	//初始化数据库
	databases.InitDB()
	//注册路由
	router.Register(r)

	//初始化定时服务
	err = service.InitCronService()
	if err != nil {
		logger.FATAL.Fatalln(err)
	}

	//初始化tg服务
	fmt.Println(configs.Get())
	if configs.Get().TGService.Enabled == true {
		tgService.InitTgService()
	}

	//启动
	err = r.Run(":8994")
	if err != nil {
		panic(err)
	}

	defer func() {
		//关闭tg服务
		if tgService.Status == 1 {
			tgService.Stop()
		}
	}()

}

package service

import (
	"errors"
	"github.com/robfig/cron/v3"
	"github.com/rroy233/EnterNEU/logger"
	"time"
)

var crontab *cron.Cron

var storageDIr = "./storage/"

func InitCronService() error {
	crontab = cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600)))

	var err error

	//清除用户过期数据
	//每天凌晨1:30点
	_, err = crontab.AddFunc("30 1 * * ?", cronClean)
	if err != nil {
		return errors.New("[定时任务][异常]添加cronClean失败:" + err.Error())
	} else {
		logger.Info.Println("[定时任务][成功]添加cronClean成功")
	}

	//更新配置文件
	//每3分钟
	_, err = crontab.AddFunc("@every 3m", cronUpdateConfig)
	if err != nil {
		return errors.New("[定时任务][异常]添加cronUpdateConfig失败:" + err.Error())
	} else {
		logger.Info.Println("[定时任务][成功]添加cronUpdateConfig成功")
	}

	crontab.Start()

	return err
}

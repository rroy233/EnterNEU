package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
	"strings"
	"time"
)

var crontab *cron.Cron

var storageDIr = "./storage/"

func InitCronService() error {
	crontab = cron.New(cron.WithLocation(time.FixedZone("CST", 8*3600)))

	var err error

	//检测e码通网址更新
	if configs.Get().General.AutoDetectUpdate == true {
		uTimes := strings.Split(configs.Get().General.AutoDetectTime, "|")
		if len(uTimes) < 1 {
			return errors.New("[定时任务][异常]添加cronDetectUpdate失败:AutoDetectTime配置项有误")
		}
		for i, uTime := range uTimes {
			u := strings.Split(uTime, ":")
			if len(u) != 2 {
				return errors.New("[定时任务][异常]添加cronDetectUpdate失败:AutoDetectTime配置项有误")
			}
			_, err = crontab.AddFunc(fmt.Sprintf("%s %s * * ?", u[1], u[0]), cronDetectUpdate)
			if err != nil {
				return errors.New(fmt.Sprintf("[定时任务][异常]添加cronDetectUpdate[%d]失败:%s", i, err.Error()))
			} else {
				logger.Info.Println(fmt.Sprintf("[定时任务][成功]添加cronDetectUpdate[%d]成功，定时%s", i, uTime))
			}
		}

	}

	//清除用户过期数据
	//每天凌晨1:30点
	_, err = crontab.AddFunc("30 1 * * ?", cronClean)
	if err != nil {
		return errors.New("[定时任务][异常]添加cronClean失败:" + err.Error())
	} else {
		logger.Info.Println("[定时任务][成功]添加cronClean成功")
	}

	//更新配置文件
	//每1分钟
	_, err = crontab.AddFunc("@every 1m", cronUpdateConfig)
	if err != nil {
		return errors.New("[定时任务][异常]添加cronUpdateConfig失败:" + err.Error())
	} else {
		logger.Info.Println("[定时任务][成功]添加cronUpdateConfig成功")
	}

	crontab.Start()

	return err
}

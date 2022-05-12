package tgService

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
)

var loggerPrefix = "[tgService]"
var serviceName = "TGService"
var Status = 0
var bot *tgbotapi.BotAPI
var cancel context.CancelFunc
var stopCtx context.Context
var cancelCh chan int
var WorkerNum int

func InitTgService() {
	if configs.Get() == nil {
		logger.FATAL.Fatalln(loggerPrefix + "配置文件未加载")
	}

	WorkerNum = configs.Get().TGService.WorkerNum

	var err error
	bot, err = tgbotapi.NewBotAPI(configs.Get().TGService.BotToken)
	if err != nil {
		logger.FATAL.Fatalln(loggerPrefix + err.Error())
	}

	//初始化许可名单
	databases.InitTGAllow()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	stopCtx, cancel = context.WithCancel(context.Background())
	cancelCh = make(chan int, WorkerNum)
	for i := 0; i < WorkerNum; i++ {
		go worker(stopCtx, updates, cancelCh)
	}

	logger.Info.Printf(loggerPrefix+"TG机器人初始化成功 %s，worker数%d", bot.Self.UserName, WorkerNum)
	Status = 1
}

func Stop() {
	cancel()
	waitForDone(cancelCh)
}

func worker(stopCtx context.Context, uc tgbotapi.UpdatesChannel, cancelCh chan int) {
	for {
		select {
		case update := <-uc:
			router(update)
		case <-stopCtx.Done():
			cancelCh <- 1
			return
		}
	}
}

func waitForDone(cancelCh chan int) {
	num := 0
	for {
		if num == WorkerNum {
			break
		}
		<-cancelCh
		num++
	}
}

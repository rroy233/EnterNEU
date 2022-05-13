package tgService

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/logger"
)

var msgQueue chan tgbotapi.Chattable

func InitSender(senderNum int) {
	msgQueue = make(chan tgbotapi.Chattable, senderNum*2)
	for i := 0; i < senderNum; i++ {
		go sender()
	}
}

func addToSendQueue(msg tgbotapi.Chattable) {
	msgQueue <- msg
}

func sender() {
	for {
		rl.Take()
		msg, _ := <-msgQueue
		_, err := bot.Send(msg)
		if err != nil {
			logger.Error.Println(loggerPrefix + "[sender]" + err.Error())
		}
	}
}

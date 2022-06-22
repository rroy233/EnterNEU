package tgService

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/configs"
)

func SendTGMsg2Admin(msg string) {
	if msg == "" {
		return
	}
	var msgConfig tgbotapi.MessageConfig
	msgConfig = tgbotapi.NewMessage(configs.Get().TGService.AdminUID, msg)
	addToSendQueue(msgConfig)
	return
}

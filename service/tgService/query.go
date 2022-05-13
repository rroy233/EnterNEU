package tgService

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"strings"
)

func queryUpload(update *tgbotapi.Update) {

	if err := callBackWithAlert(update.CallbackQuery.ID, "请使用网页版上传"); err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
	}

	return
}

func queryDel(update *tgbotapi.Update) {
	if err := callBack(update.CallbackQuery.ID, "ok"); err != nil {
		logger.Error.Println(loggerPrefix + getLogPrefixCallbackQuery(update) + err.Error())
	}
	token := strings.Split(update.CallbackQuery.Data, "#")[1]

	helper := databases.NewHelper("")
	helper.SetToken(token)
	if helper.Validate() == false {
		if _, err := bot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "过期已删除,可使用 /new 重新创建")); err != nil {
			logger.Error.Println(loggerPrefix+getLogPrefixCallbackQuery(update)+"[EditMessageText]", err)
		}
	} else {
		if err := helper.Delete(); err != nil {
			logger.Error.Println(loggerPrefix + getLogPrefixCallbackQuery(update) + err.Error())
			err = callBackWithAlert(update.CallbackQuery.ID, "删除失败")
		}
		if _, err := bot.Send(tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "已删除,可使用 /new 重新创建")); err != nil {
			logger.Error.Println(loggerPrefix+getLogPrefixCallbackQuery(update)+"[EditMessageText]", err)
		}
	}

	logger.Info.Printf(loggerPrefix + getLogPrefixCallbackQuery(update) + "删除成功")

	return
}

func queryGetVideo(update *tgbotapi.Update) {
	if len(strings.Split(update.CallbackQuery.Data, "#")) != 3 {
		sendPlainText(update, "参数无效")
		return
	}
	token := strings.Split(update.CallbackQuery.Data, "#")[1]
	key := strings.Split(update.CallbackQuery.Data, "#")[2]

	//修改原信息
	edit := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, getTicketInlineKeyboardNoVideo(token, key))
	addToSendQueue(edit)

	//callback
	err := callBack(update.CallbackQuery.ID, "ok")
	if err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
	}

	//发送
	text := "Shadowrocket配置视频教程:\n点我"
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	msg.Entities = []tgbotapi.MessageEntity{
		entityTextLink(text, "点我", "https://upload.enterneu.icu/enterNEU_ios_compressed.mp4"),
	}
	addToSendQueue(msg)
	return
}

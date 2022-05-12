package tgService

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
)

func getCodeTypeKeyboard() (*tgbotapi.ReplyKeyboardMarkup, error) {
	//每行四个
	kb := new(tgbotapi.ReplyKeyboardMarkup)
	kb.ResizeKeyboard = true
	codeConsts, err := configs.GetECodeConst()
	if err != nil {
		return nil, err
	}
	bts := make([][]tgbotapi.KeyboardButton, len(codeConsts.CodeTypes)%4)
	for i := 0; i < len(codeConsts.CodeTypes)%4; i++ {
		bts[i] = make([]tgbotapi.KeyboardButton, 0)
		for j := 0; j < 4; j++ {
			if (i*4 + j) < len(codeConsts.CodeTypes) {
				bts[i] = append(bts[i], tgbotapi.KeyboardButton{
					Text: codeConsts.CodeTypes[i*4+j],
				})
			}
		}
	}
	kb.Keyboard = bts
	return kb, err
}

func getActualVehicleKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	//每行2个
	kb := new(tgbotapi.ReplyKeyboardMarkup)
	kb.ResizeKeyboard = true
	kb.Keyboard = [][]tgbotapi.KeyboardButton{
		{tgbotapi.KeyboardButton{
			Text: "入",
		}, tgbotapi.KeyboardButton{
			Text: "出",
		}},
	}
	return kb
}
func getExpTimeKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	//每行2个
	kb := new(tgbotapi.ReplyKeyboardMarkup)
	kb.ResizeKeyboard = true
	kb.Keyboard = [][]tgbotapi.KeyboardButton{
		{tgbotapi.KeyboardButton{
			Text: "1小时",
		}, tgbotapi.KeyboardButton{
			Text: "24小时",
		}},
		{tgbotapi.KeyboardButton{
			Text: "3天",
		}, tgbotapi.KeyboardButton{
			Text: "一周",
		}},
	}
	return kb
}

func getUploadKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	//每行2个
	kb := new(tgbotapi.ReplyKeyboardMarkup)
	kb.ResizeKeyboard = true
	kb.Keyboard = [][]tgbotapi.KeyboardButton{
		{tgbotapi.KeyboardButton{
			Text: "上传头像",
		}, tgbotapi.KeyboardButton{
			Text: "不上传",
		}},
	}
	return kb
}

func clearKeyboard(update *tgbotapi.Update, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := bot.Send(msg)
	if err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
	}
}

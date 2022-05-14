package tgService

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"strings"
)

const noPermissionText = "您没有使用的权限!\n如果您是东大的学生，您可以向已被授权的东大学生索要授权。"

func router(update tgbotapi.Update) {
	//debug
	if configs.Get().General.Production == false {
		jsonlog, _ := json.Marshal(update)
		logger.Info.Println(loggerPrefix + "[update]" + string(jsonlog))
	}
	//debug

	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	//拦截无权限访问
	if update.Message != nil && databases.GetTGAllow().IsIn(update.Message.From.ID) == false {
		sendSticker(&update, "CAACAgUAAxkBAAMHYnsNDA4Eq1FyWvpoCCZhhud7hE8AAp8AA9x4DAiSZEHHuSz0WCQE")
		text := fmt.Sprintf("%s\n这是您的UID:%d\n如有疑问请加入东北大学群组", noPermissionText, update.Message.From.ID)
		sendPlainText(&update, text,
			entityTextLink(text, "东北大学群组", "https://t.me/dongbeidaxue"),
			entityBold(text, "东大的学生"),
		)
		return
	}
	if update.CallbackQuery != nil && databases.GetTGAllow().IsIn(update.CallbackQuery.From.ID) == false {
		if err := callBackWithAlert(update.CallbackQuery.ID, noPermissionText); err != nil {
			logger.Error.Println(loggerPrefix + err.Error())
		}
		return
	}

	//判断是否为群组
	if update.Message != nil && update.Message.Chat.IsGroup() == true {
		sendPlainText(&update, "暂时不支持群组")
		return
	}

	//判断是否包含指令
	if update.Message != nil && update.Message.IsCommand() == true {
		//go on
		command := update.Message.Command()
		logger.Info.Printf("%s[Message][User:%d @%s %s][Chat:%s]Command:%s",
			loggerPrefix,
			update.Message.From.ID,
			update.Message.From.UserName,
			update.Message.From.FirstName+update.Message.From.LastName,
			fmt.Sprintf("(%s) %d %s", update.Message.Chat.Type, update.Message.Chat.ID, update.Message.Chat.Title),
			update.Message.Text,
		)
		if command != "" {
			switch command {
			case "start":
				commandStart(&update)
			case "help":
				commandHelp(&update)
			case "new":
				commandNew(&update)
			case "cancel":
				commandCancel(&update)
			case "clear":
				commandCleanKeyBoard(&update)
			case "add":
				commandAddPermission(&update)
			case "remove":
				commandRemovePermission(&update)
			case "list":
				commandListPermission(&update)
			default:
				sendPlainText(&update, "命令不存在")
				return
			}
		}
		return
	}

	//判断是否为callback_query
	if update.CallbackQuery != nil {
		logger.Info.Printf("%s[CallbackQuery][User:%d @%s %s][Chat:%s]Data:%s",
			loggerPrefix,
			update.CallbackQuery.From.ID,
			update.CallbackQuery.From.UserName,
			update.CallbackQuery.From.FirstName+update.CallbackQuery.From.LastName,
			fmt.Sprintf("(%s) %d %s", update.CallbackQuery.Message.Chat.Type, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Chat.Title),
			update.CallbackQuery.Data,
		)
		//go on
		switch strings.Split(update.CallbackQuery.Data, "#")[0] {
		case "upload":
			queryUpload(&update)
		case "del":
			queryDel(&update)
		case "get_video":
			queryGetVideo(&update)
		default:
			if err := callBackWithAlert(update.CallbackQuery.ID, "操作不存在"); err != nil {
				logger.Error.Println(loggerPrefix+"[route]default发送CallBackWithAlert失败", err)
			}
		}
		return
	}

	//处理表单
	formHandler(&update)

	//忽略
	return
}

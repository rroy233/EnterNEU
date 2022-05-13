package tgService

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"strconv"
	"strings"
)

const commandList = `
/help - 查看帮助[私聊]
/new - 开始创建凭证[私聊]
/cancel - 取消创建凭证[私聊]
/clear - 清除小键盘[私聊]
/add <uid> - 给予用户使用权限[私聊]
/remove <uid>|all - 撤销用户使用权限[私聊][管理员]
`

func commandStart(update *tgbotapi.Update) {
	sendPlainText(update, "欢迎使用\n使用 /help 查看帮助")
	return
}
func commandHelp(update *tgbotapi.Update) {
	sendPlainText(update, commandList)
	return
}

func commandCancel(update *tgbotapi.Update) {
	session, _ := getFormSession(update.Message.From.ID)
	if session != nil {
		err := session.del()
		if err != nil {
			sendPlainText(update, "删除失败")
			return
		}
	} else {
		sendPlainText(update, "您没有活跃中的会话")
		return
	}
	clearKeyboard(update, "已结束创建。")
	return
}

func commandCleanKeyBoard(update *tgbotapi.Update) {
	clearKeyboard(update, "已执行")
	return
}

func commandNew(update *tgbotapi.Update) {
	session := newFormSession(update.Message.From.ID, stepStarted)
	err := session.save()
	if err != nil {
		sendPlainText(update, "系统异常(保存session失败)")
		return
	}

	keybord := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{
		{"开始", false, false, nil},
		{"/cancel", false, false, nil},
	})
	text := "现在开始创建您的凭证:\n说明文档:点击查看 \n\n(创建过程中可点击 /cancel 取消)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.Entities = []tgbotapi.MessageEntity{
		entityTextLink(text, "点击查看", "https://telegra.ph/EnterNEU%E8%AF%B4%E6%98%8E%E6%96%87%E6%A1%A3-05-10"),
	}
	msg.ReplyMarkup = keybord
	addToSendQueue(msg)
	return
}

func commandAddPermission(update *tgbotapi.Update) {
	if update.Message.Text == "/add" || len(strings.Split(update.Message.Text, " ")) < 2 {
		text := "该功能允许所有已授权的用户使用\n【授权原则】\n仅允许授权给东大学生!!!\n请务必严格验证其身份。\n\n请按照下面的格式进行发送:\n /add <UID> (如/add 123)"
		sendPlainText(update, text,
			entityBold(text, "仅允许授权给东大学生"),
			entityUnderline(text, "严格验证其身份"),
		)
		return
	}

	playLoad := strings.Split(update.Message.Text, " ")[1]
	uid, err := strconv.ParseInt(playLoad, 10, 64)
	if err != nil {
		logger.Error.Println(loggerPrefix + getLogPrefixMessage(update) + "UID解析失败:" + err.Error())
		sendPlainText(update, "UID解析失败")
		return
	}

	chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
		tgbotapi.ChatConfig{
			ChatID: uid,
		},
	})
	if err != nil {
		logger.Error.Println(loggerPrefix + getLogPrefixMessage(update) + "校验UID失败:" + err.Error())
		sendPlainText(update, "UID无效")
		return
	}
	if chat.IsPrivate() == false {
		sendPlainText(update, "不允许添加群组!!!!")
		return
	}
	logger.Info.Printf("%s%s已添加用户:%d@%s", loggerPrefix, getLogPrefixMessage(update), chat.ID, chat.UserName)
	logger.Debug.Println(chat)

	//给该用户推送通知
	msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf(
		"用户%s(@%s)已批准您使用本bot\n\n您也可以使用 /add 命令批准别的用户。\n注意,请严格按照授权原则进行使用。",
		update.Message.From.FirstName+update.Message.From.LastName,
		update.Message.From.UserName))
	addToSendQueue(msg)

	//加入许可名单
	databases.GetTGAllow().Put(uid)

	//回复操作者
	sendPlainText(update, fmt.Sprintf("添加UID:%d成功", uid))
	return
}

func commandRemovePermission(update *tgbotapi.Update) {
	//判断是否为管理员
	if update.Message.From.ID != databases.GetTGAllow().AdminUID {
		sendPlainText(update, "本指令仅允许管理员使用")
		return
	}
	if update.Message.Text == "/remove" || len(strings.Split(update.Message.Text, " ")) < 2 {
		sendPlainText(update, "格式:\n /remove <UID> (如/remove 123)\n删除所有用户则输入 /remove all")
		return
	}

	playLoad := strings.Split(update.Message.Text, " ")[1]
	if playLoad == "all" {
		databases.GetTGAllow().FlushAll()
		sendPlainText(update, "操作成功")
		return
	} else {
		uid, err := strconv.ParseInt(playLoad, 10, 64)
		if err != nil {
			logger.Error.Println(loggerPrefix + getLogPrefixMessage(update) + "UID解析失败:" + err.Error())
			sendPlainText(update, "UID解析失败")
			return
		}

		//判断是否在名单上
		if databases.GetTGAllow().IsIn(uid) == false {
			sendPlainText(update, "此UID不在许可名单上")
			return
		}

		chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{
			tgbotapi.ChatConfig{
				ChatID: uid,
			},
		})
		if err != nil {
			logger.Error.Println(loggerPrefix + getLogPrefixMessage(update) + "校验UID失败:" + err.Error())
			sendPlainText(update, "UID无效")
			return
		}
		if chat.IsPrivate() == false {
			sendPlainText(update, "不操作群组!!!!")
			return
		}
		//加入许可名单
		databases.GetTGAllow().Remove(uid)

		logger.Info.Printf("%s%s已移除用户:%d@%s", loggerPrefix, getLogPrefixMessage(update), chat.ID, chat.UserName)
		//给该用户推送通知
		msg := tgbotapi.NewMessage(chat.ID, fmt.Sprintf(
			"用户%s(@%s)已移除您使用本bot的权限",
			update.Message.From.FirstName+update.Message.From.LastName,
			update.Message.From.UserName))
		addToSendQueue(msg)
		//回复操作者
		sendPlainText(update, fmt.Sprintf("移除UID:%d成功", uid))
	}

	return
}

func commandListPermission(update *tgbotapi.Update) {
	//判断是否为管理员
	if update.Message.From.ID != databases.GetTGAllow().AdminUID {
		sendPlainText(update, "本指令仅允许管理员使用")
		return
	}

	list := databases.GetTGAllow().List()
	//Message.text最多4096个字符
	text := "【当前已授权的名单】\n" + list
	if len([]rune(text)) > 4096 {
		text = string([]rune(text)[:4090]) + "..."
	}
	sendPlainText(update, text)
	return
}

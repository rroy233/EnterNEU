package tgService

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
	"time"
)

func formHandler(update *tgbotapi.Update) {
	session, err := getFormSession(update.Message.From.ID)
	if err != nil {
		sendPlainText(update, "请使用 /new 创建新凭证")
		return
	}

	//判断是否发送图片
	text := update.Message.Text
	if session.Step == stepFormDecideUpload {
		if text == "上传头像" {
			clearKeyboard(update, "请以图片方式发送您的头像:(大小不超过1MB)")
			_ = session.setStep(stepFormUpload).save()
			return
		} else {
			session.Step = stepFormFinished
		}
	}
	//处理发送的图片
	if session.Step == stepFormUpload {
		if update.Message.Photo != nil && len(update.Message.Photo) != 0 {
			err = session.storeImg(update)
			if err != nil {
				sendPlainText(update, fmt.Sprintf("失败，请重新上传(%s):", err.Error()))
				return
			}
			session.setStep(stepFormFinished)
		} else {
			sendPlainText(update, "请以图片方式发送您的头像:(大小不超过1MB)")
			return
		}
	}

	switch session.Step {
	case stepStarted:
		session.Form.Key = utils.NewUUIDToken()
		session.Form.KeyMD5 = utils.MD5Short(session.Form.Key + configs.Get().General.Md5Salt)
		session.setStep(stepFormName)
		clearKeyboard(update, "开始")
		_ = session.save()
		sendPlainText(update, "已为您随机生成秘钥:\n"+session.Form.Key)
		time.Sleep(200 * time.Millisecond)
		hint := "请输入您的姓名(格式:孙*川):"
		sendPlainText(update, hint, entityBold(hint, "姓名"))
	case stepFormName:
		session.Form.Name = text
		session.setStep(stepFormStuID)
		_ = session.save()
		hint := "请输入您的学号:"
		sendPlainText(update, hint, entityBold(hint, "学号"))
	case stepFormStuID:
		session.Form.StuID = text
		session.setStep(stepFormEntranceName)
		_ = session.save()
		hint := "请输入入口名(格式:xx校区xx门):"
		sendPlainText(update, hint, entityBold(hint, "入口名"))
	case stepFormEntranceName:
		session.Form.EntranceName = text
		session.setStep(stepFormCodeType)
		_ = session.save()
		kb, _ := getCodeTypeKeyboard()
		hint := "请选择(小键盘上的)提示文本:"
		sendPlainTextWithKeyboard(update, hint, kb, entityBold(hint, "提示文本"))
	case stepFormCodeType:
		ecc, err := configs.GetECodeConst()
		if err != nil {
			logger.Error.Println(loggerPrefix + "读取e-code配置文件失败:" + err.Error())
			sendPlainText(update, "读取e-code配置文件失败")
			return
		}
		if ecc.CodeTypeIndexByText[text] == "" {
			sendPlainText(update, "提示文本无效,请重新选择", entityBold("提示文本无效,请重新选择", "提示文本"))
			return
		}
		session.Form.CodeType = ecc.CodeTypeIndexByText[text]
		session.setStep(stepFormActualVehicle)
		_ = session.save()
		hint := "请选择(小键盘上的)方向:"
		sendPlainTextWithKeyboard(update, hint, getActualVehicleKeyboard(), entityBold(hint, "方向"))
	case stepFormActualVehicle:
		if text == "出" {
			session.Form.ActualVehicle = "1"
		} else if text == "入" {
			session.Form.ActualVehicle = "0"
		} else {
			sendPlainText(update, "输入无效")
			return
		}
		session.setStep(stepFormExpTimeType)
		_ = session.save()

		hint := "请选择(小键盘上的)过期时间:"
		sendPlainTextWithKeyboard(update, hint, getExpTimeKeyboard(), entityBold(hint, "过期时间"))
		//调出小键盘
	case stepFormExpTimeType:
		timeType := 0
		switch text {
		case "1小时":
			timeType = 0
		case "24小时":
			timeType = 1
		case "3天":
			timeType = 2
		case "一周":
			timeType = 3
		default:
			sendPlainText(update, "输入无效")
			return
		}
		session.Form.ExpTimeType = timeType
		session.setStep(stepFormDecideUpload)
		_ = session.save()
		sendPlainTextWithKeyboard(update, "请选择(小键盘上的)是否上传头像:\n(若此时不上传，之后请使用网页版上传)", getUploadKeyboard())
		//调出小键盘
	case stepFormFinished:
		token, err := session.createEcode()
		if err != nil {
			sendPlainText(update, "异常")
			return
		}
		clearKeyboard(update, "已完成!")
		session.sendTicketMsg(update, token)
	default:
		sendPlainText(update, "当前状态无效，请使用 /cancel 取消。")
		return
	}
	return
}

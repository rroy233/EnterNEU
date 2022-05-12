package tgService

import (
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type sessionStep int

var sessionExpTime = time.Minute * 5

//一次创建凭证的过程中，用户当前处在的状态
const (
	stepStarted = sessionStep(iota)
	stepFormName
	stepFormStuID
	stepFormEntranceName
	stepFormCodeType
	stepFormActualVehicle
	stepFormExpTimeType
	stepFormDecideUpload
	stepFormUpload
	stepFormFinished
)

const (
	expTimeTypeOneHour = iota
	expTimeTypeOneDay
	expTimeTypeThreeDays
	expTimeTypeOneWeek
)

type tgFormSession struct {
	SessionID   string      `json:"sessionID"`
	UserID      int64       `json:"userID"`
	Step        sessionStep `json:"step"`
	TicketMsgID int         `json:"ticketMsgID"`
	Token       string      `json:"token"`
	Form        struct {
		Key           string `json:"key"`
		KeyMD5        string `json:"keyMD5"`
		Name          string `json:"name"`
		StuID         string `json:"stuID"`
		EntranceName  string `json:"entranceName"`
		CodeType      string `json:"codeType"`
		ActualVehicle string `json:"actualVehicle"`
		ExpTimeType   int    `json:"expTimeType"`
	} `json:"form"`
	ImgUploaded    bool   `json:"ImgUploaded"`
	TmpFilePath    string `json:"tmpFilePath"`
	ImgContentType string `json:"ImgContentType"`
	ExpTime        string `json:"expTime"`
}

func getFormSession(UID int64) (*tgFormSession, error) {
	data, err := databases.GetFromRedis(serviceName, []string{
		"FormSession",
		fmt.Sprintf("USER_%d", UID),
	})
	if err != nil {
		return nil, err
	}

	fs := new(tgFormSession)
	err = json.Unmarshal([]byte(data), fs)
	if err != nil {
		return nil, err
	}
	return fs, nil
}

func newFormSession(UID int64, step sessionStep) *tgFormSession {
	f := new(tgFormSession)
	f.UserID = UID
	f.Step = step
	f.SessionID = utils.NewUUIDToken()
	return f
}

func (fs *tgFormSession) setStep(step sessionStep) *tgFormSession {
	fs.Step = step
	return fs
}

func (fs tgFormSession) save() error {
	data, err := json.Marshal(fs)
	if err != nil {
		return err
	}

	err = databases.SaveToRedis(serviceName, []string{
		"FormSession",
		fmt.Sprintf("USER_%d", fs.UserID),
	}, data, sessionExpTime)
	return err
}

func (fs tgFormSession) Update(duration time.Duration) error {
	data, err := json.Marshal(fs)
	if err != nil {
		return err
	}

	err = databases.SaveToRedis(serviceName, []string{
		"FormSession",
		fmt.Sprintf("USER_%d", fs.UserID),
	}, data, duration)
	return err
}

func (fs tgFormSession) del() error {
	err := databases.DeleteFromRedis(serviceName, []string{
		"FormSession",
		fmt.Sprintf("USER_%d", fs.UserID),
	})
	return err
}

func (fs *tgFormSession) createEcode() (string, error) {
	if fs.Step != stepFormFinished {
		return "", errors.New("未完成表单")
	}
	form := fs.Form

	//判断key是否有效
	if form.KeyMD5 != utils.MD5Short(form.Key+configs.Get().General.Md5Salt) {
		return "", errors.New("key无效")
	}

	//判断expDuration
	expDuration := 50 * time.Minute
	switch form.ExpTimeType {
	case expTimeTypeOneHour:
		expDuration = time.Hour
	case expTimeTypeOneDay:
		expDuration = 24 * time.Hour
	case expTimeTypeThreeDays:
		expDuration = 3 * 24 * time.Hour
	case expTimeTypeOneWeek:
		expDuration = 7 * 24 * time.Hour
	default:
		return "", errors.New("ExpTimeType无效")
	}

	fs.ExpTime = time.Now().Add(expDuration).Format("2006-01-02 15:04:05")

	//存储
	helper := databases.NewHelper(form.Key)
	token, err := helper.CreateECode(form.Name, form.StuID, form.EntranceName, form.Key, "", form.CodeType, form.ActualVehicle, expDuration)
	if err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
		return "", errors.New("创建失败")
	}

	//处理图片,从临时目录移动出来
	if fs.ImgUploaded == true {
		tempFile, err := ioutil.ReadFile(fs.TmpFilePath)
		if err != nil {
			logger.Error.Println(loggerPrefix + err.Error())
			return "", errors.New("找不到临时文件")
		}

		//加密
		imgEncrypt, err := utils.AesCbcEncrypt(tempFile, fs.Form.Key)
		if err = ioutil.WriteFile(fmt.Sprintf("./storage/%s.data", token), imgEncrypt, 0755); err != nil {
			logger.Error.Println(loggerPrefix + err.Error())
			return "", errors.New("存储图片发生错误")
		}

		//更新数据库
		if err = helper.UpdateImg(fmt.Sprintf("./storage/%s.data", token), fs.ImgContentType); err != nil {
			logger.Error.Println(loggerPrefix + err.Error())
			return "", errors.New("更新数据库发生错误")
		}

		//移除临时文件
		if err = os.Remove(fs.TmpFilePath); err != nil {
			logger.Error.Println(loggerPrefix + err.Error())
		}
	}

	//删除session
	err = fs.del()
	if err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
	}

	return token, err
}

func (fs *tgFormSession) storeImg(update *tgbotapi.Update) error {
	msg, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "处理中。。。"))
	if err != nil {
		logger.Error.Println(loggerPrefix + "[storeImg]" + err.Error())
	}

	imgIndex := len(update.Message.Photo) - 1 //原图
	remoteFile, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: update.Message.Photo[imgIndex].FileID,
	})
	if err != nil {
		log.Println("[util][upload.UploadCommand]获取文件失败:", err)
		return errors.New("获取文件失败(-1)")
	}
	filePath, err := downloadFile(remoteFile.Link(bot.Token))
	if err != nil {
		log.Println("[util][upload.UploadCommand]下载文件失败:", err)
		return errors.New("获取文件失败(-2)")
	}

	fs.ImgUploaded = true
	fs.TmpFilePath = filePath
	fs.ImgContentType = "image/jpeg"

	//撤回处理进度消息
	if _, err := bot.Request(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, msg.MessageID)); err != nil {
		logger.Error.Println(loggerPrefix + "[storeImg]" + err.Error())
	}
	return nil
}

func (fs *tgFormSession) sendTicketMsg(update *tgbotapi.Update, token string) {

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("e码通页面", fmt.Sprintf("%s/%s/%s", configs.Get().General.BaseUrl, token, fs.Form.Key)),
			tgbotapi.NewInlineKeyboardButtonURL("管理面板(web)", fmt.Sprintf("%s/#/status/%s/%s", configs.Get().General.BaseUrl, token, fs.Form.Key)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("删除", fmt.Sprintf("del#%s#%s", token, fs.Form.Key)),
		),
	)

	text := fmt.Sprintf("%s\n%s\n%s\n%s",
		"#EnterNEU凭证\n",
		"【token】"+token+"\n【key】"+fs.Form.Key+"\n【过期时间】"+fs.ExpTime+"\n",
		"e码通页面入口:\n"+fmt.Sprintf("%s/%s/%s", configs.Get().General.BaseUrl, token, fs.Form.Key),
		"Shadowrocket配置地址:\n"+fmt.Sprintf("%s/%s/%s/shadowrocket", configs.Get().General.BaseUrl, token, fs.Form.Key),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboard
	msg.Entities = []tgbotapi.MessageEntity{
		entityTag(text, "#EnterNEU凭证"),
	}
	_, err := bot.Send(msg)
	if err != nil {
		logger.Error.Println(loggerPrefix + err.Error())
	}

	return
}

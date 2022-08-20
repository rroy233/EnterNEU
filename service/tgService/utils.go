package tgService

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rroy233/EnterNEU/logger"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf16"
)

func sendPlainText(update *tgbotapi.Update, text string, entity ...tgbotapi.MessageEntity) {
	if update.Message == nil {
		return
	}
	var msg tgbotapi.MessageConfig
	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyToMessageID = update.Message.MessageID
		if entity != nil {
			msg.Entities = entity
		}
		addToSendQueue(msg)
	} else if update.CallbackQuery != nil || update.CallbackQuery.Message != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		if entity != nil {
			msg.Entities = entity
		}
		addToSendQueue(msg)
	}
}

func sendImg(update *tgbotapi.Update, fileData []byte) (msgSent tgbotapi.Message) {
	if update.Message == nil {
		return
	}
	var msg tgbotapi.PhotoConfig
	file := tgbotapi.FileBytes{
		Name:  "image.jpg",
		Bytes: fileData,
	}
	if update.Message != nil {
		msg = tgbotapi.NewPhoto(update.Message.Chat.ID, file)
	} else if update.CallbackQuery != nil || update.CallbackQuery.Message != nil {
		msg = tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, file)
	}
	smsg, _ := bot.Send(msg)
	return smsg
}

func sendSticker(update *tgbotapi.Update, fileID string) {
	if update.Message == nil {
		return
	}
	var msg tgbotapi.StickerConfig
	if update.Message != nil {
		msg = tgbotapi.NewSticker(update.Message.Chat.ID, tgbotapi.FileID(fileID))
		msg.ReplyToMessageID = update.Message.MessageID
		addToSendQueue(msg)
	} else if update.CallbackQuery != nil || update.CallbackQuery.Message != nil {
		msg = tgbotapi.NewSticker(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileID(fileID))
		addToSendQueue(msg)
	}
}

func sendPlainTextWithKeyboard(update *tgbotapi.Update, text string, keyboard *tgbotapi.ReplyKeyboardMarkup, entity ...tgbotapi.MessageEntity) {
	if update.Message == nil {
		return
	}
	var msg tgbotapi.MessageConfig
	if update.Message != nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = *keyboard
		if entity != nil {
			msg.Entities = entity
		}
		addToSendQueue(msg)
	} else if update.CallbackQuery != nil || update.CallbackQuery.Message != nil {
		msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		msg.ReplyMarkup = *keyboard
		if entity != nil {
			msg.Entities = entity
		}
		addToSendQueue(msg)
	}
}

func downloadFile(fileUrl string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fileUrl, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//get file name
	oFileName := "file"
	urls := strings.Split(fileUrl, "/")
	if len(urls) == 0 {
		return "", errors.New("url无效")
	}
	if strings.Contains(urls[len(urls)-1], ".") != false {
		oFileName = urls[len(urls)-1]
	}

	fileName := fmt.Sprintf("./storage/tmp/upload_%d_%s", time.Now().UnixMicro(), oFileName)
	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func callBack(callbackQueryID string, text string) error {
	callback := tgbotapi.NewCallback(callbackQueryID, text)
	//不能用bot.Send(callback)方法，有bug
	resp, err := bot.Request(callback)
	if err != nil {
		return err
	}
	if string(resp.Result) != "true" {
		return errors.New("请求不ok")
	}
	return err
}
func callBackWithAlert(callbackQueryID string, text string) error {
	callback := tgbotapi.NewCallbackWithAlert(callbackQueryID, text)
	//不能用bot.Send(callback)方法，有bug
	resp, err := bot.Request(callback)
	if err != nil {
		return err
	}
	if string(resp.Result) != "true" {
		return errors.New("请求不ok")
	}

	return err
}

func getLogPrefixMessage(update *tgbotapi.Update) string {
	return fmt.Sprintf("[Message][User:%d @%s %s][Chat:%s]",
		update.Message.From.ID,
		update.Message.From.UserName,
		update.Message.From.FirstName+update.Message.From.LastName,
		fmt.Sprintf("(%s) %d %s", update.Message.Chat.Type, update.Message.Chat.ID, update.Message.Chat.Title),
	)
}

func getLogPrefixCallbackQuery(update *tgbotapi.Update) string {
	return fmt.Sprintf("[CallbackQuery][User:%d @%s %s][Chat:%s]",
		update.CallbackQuery.From.ID,
		update.CallbackQuery.From.UserName,
		update.CallbackQuery.From.FirstName+update.CallbackQuery.From.LastName,
		fmt.Sprintf("(%s) %d %s", update.CallbackQuery.Message.Chat.Type, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.Chat.Title),
	)
}

//顺序查找UTF-6编码字符串中子串的第一次出现的位置
//返回offset=-1则为找到
func getPartIndex(text, part string) (offset, length int) {
	textUTF16 := utf16.Encode([]rune(text))
	PartUTF16 := utf16.Encode([]rune(part))
	offset = 0
	i := 0
	j := 0

	//debug
	//log.Println("textUTF16=", textUTF16)
	//log.Println("PartUTF16=", PartUTF16)

	for {
		//越界
		if i > len(textUTF16)-1 || j > len(PartUTF16) {
			offset = -1
			break
		}
		//debug
		//log.Printf("offset=%d,textUTF16[%d]=%d,PartUTF16[%d]=%d\n", offset, i, textUTF16[i], j, PartUTF16[j])

		//判断
		if textUTF16[i] == PartUTF16[j] {
			i++
			j++
		} else {
			j = 0
			offset++
			i = offset
		}

		//结果
		if j > len(PartUTF16)-1 {
			break
		}
	}
	if offset == -1 {
		if logger.Error != nil {
			logger.FATAL.Printf("%s[util][getPartIndex]无法在【%s】中找到【%s】\n", loggerPrefix, text, part)
		} else {
			log.Printf("%s[util][getPartIndex]无法在【%s】中找到【%s】\n", loggerPrefix, text, part)
		}
	}

	return offset, len(PartUTF16)
}

func entityBold(text, boldPart string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, boldPart)
	return tgbotapi.MessageEntity{
		Type:   "bold",
		Offset: offset,
		Length: length,
	}
}

func entityUnderline(text, boldPart string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, boldPart)
	return tgbotapi.MessageEntity{
		Type:   "underline",
		Offset: offset,
		Length: length,
	}
}

func entityLink(text, part, url string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, part)
	return tgbotapi.MessageEntity{
		Type:   "url",
		Offset: offset,
		Length: length,
		URL:    url,
	}
}

func entityMention(text, part string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, part)
	return tgbotapi.MessageEntity{
		Type:   "mention",
		Offset: offset,
		Length: length,
	}
}

func entityTag(text, part string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, part)
	return tgbotapi.MessageEntity{
		Type:   "hashtag",
		Offset: offset,
		Length: length,
	}
}

func entityCode(text, part string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, part)
	return tgbotapi.MessageEntity{
		Type:   "code",
		Offset: offset,
		Length: length,
	}
}

func entityTextLink(text, part, url string) tgbotapi.MessageEntity {
	offset, length := getPartIndex(text, part)
	return tgbotapi.MessageEntity{
		Type:   "text_link",
		Offset: offset,
		Length: length,
		URL:    url,
	}
}

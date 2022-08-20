package captcha

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/steambap/captcha"
	"time"
)

type CaptchaRes struct {
	Img     []byte
	UUID    string
	Text    string
	ExpTIme int64
	MsgID   int
}

// Create
// 生成验证码
func Create() (cr CaptchaRes, err error) {
	loggerPrefix := "[HandleCaptcha]"

	data, err := captcha.New(210, 70)
	if err != nil {
		logger.Error.Println(loggerPrefix+"生成验证码失败：", err)
		return CaptchaRes{}, errors.New("生成验证码失败")
	}

	img := bytes.NewBuffer(nil)
	err = data.WriteImage(img)
	if err != nil {
		logger.Error.Println(loggerPrefix+"验证码图片输出失败：", err)
		return CaptchaRes{}, errors.New("验证码图片输出失败")
	}

	res := CaptchaRes{
		Img:     img.Bytes(),
		UUID:    uuid.New().String(),
		Text:    data.Text,
		ExpTIme: time.Now().Add(5 * time.Minute).Unix(),
	}

	return res, nil
}

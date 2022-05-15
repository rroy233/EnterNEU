package databases

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
	"os"
	"time"
)

type helper struct {
	init  bool
	token string
	key   string
}

type StoreStruct struct {
	Token   string `json:"Token"`
	Student struct {
		ImgUploaded    bool   `json:"ImgUploaded"`
		ImgPath        string `json:"ImgPath"`
		ImgContentType string `json:"ImgContentType"`
		ID             string `json:"ID"`
		Name           string `json:"Name"`
	} `json:"Student"`
	Entrance struct {
		Name          string `json:"Name"`
		CodeType      string `json:"CodeType"`
		ActualVehicle string `json:"actualVehicle"`
	}
	ExpTime int64 `json:"ExpTime"`
}

func NewHelper(key string) *helper {
	th := new(helper)
	th.key = key
	th.init = true
	return th
}

func (th *helper) SetToken(token string) *helper {
	th.token = token
	return th
}

func (th helper) Validate() bool {
	if th.token == "" {
		return false
	}
	ex := rdb.Exists(utils.Ctx(), fmt.Sprintf("EnterNeu:Ticket:%s", th.token)).Val()
	if ex == 1 {
		return true
	}
	return false
}

func (th helper) Get() (*StoreStruct, error) {
	data := rdb.Get(utils.Ctx(), fmt.Sprintf("EnterNeu:Ticket:%s", th.token)).Val()
	if data == "" {
		return nil, errors.New(fmt.Sprintf("token:%s不存在", th.token))
	}

	store := new(StoreStruct)
	err := json.Unmarshal([]byte(data), store)
	return store, err
}

// CreateECode 创建e-code
//
// 返回token
func (th *helper) CreateECode(stuName, stuID, Entrance, key, imgPath, codeType, actualVehicle string, expTime time.Duration) (token string, err error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	token = utils.MD5Short(uid.String())
	th.token = token

	store := new(StoreStruct)
	store.Token = token
	if imgPath != "" {
		store.Student.ImgUploaded = true
		store.Student.ImgPath = imgPath
	}
	store.Student.ID, err = utils.AesCbcEncryptString(stuID, key)
	store.Student.Name, err = utils.AesCbcEncryptString(stuName, key)
	store.Entrance.Name = Entrance
	store.ExpTime = time.Now().Add(expTime).Unix()
	store.Entrance.CodeType = codeType
	store.Entrance.ActualVehicle = actualVehicle

	storeJson, err := json.Marshal(store)
	if err != nil {
		return "", err
	}
	err = rdb.Set(utils.Ctx(), fmt.Sprintf("EnterNeu:Ticket:%s", token), string(storeJson), expTime).Err()
	return token, err
}

// UpdateImg 更新图片地址
func (th helper) UpdateImg(imgPath, contentType string) (err error) {
	if th.token == "" {
		return errors.New("无token")
	}

	od, err := rdb.Get(utils.Ctx(), fmt.Sprintf("EnterNeu:Ticket:%s", th.token)).Result()
	if err != nil {
		return err
	}
	old := new(StoreStruct)
	err = json.Unmarshal([]byte(od), old)
	if err != nil {
		return err
	}

	if imgPath != "" {
		old.Student.ImgUploaded = true
		old.Student.ImgPath = imgPath
		old.Student.ImgContentType = contentType
	} else {
		if old.Student.ImgUploaded == true {
			//改为空
			old.Student.ImgUploaded = false
			old.Student.ImgPath = ""
			old.Student.ImgContentType = ""
		}
	}

	storeJson, err := json.Marshal(old)
	if err != nil {
		return err
	}
	err = rdb.Set(utils.Ctx(), fmt.Sprintf("EnterNeu:Ticket:%s", th.token), string(storeJson), redis.KeepTTL).Err()
	return err
}

func (th helper) Delete() error {
	store, err := th.Get()
	if err != nil {
		return errors.New("记录不存在")
	}

	//判断头像是否存在
	if store.Student.ImgUploaded == true {
		if imgExist, err := utils.FsExists(store.Student.ImgPath); err == nil && imgExist == true {
			_ = os.Remove(store.Student.ImgPath)
		} else {
			logger.Error.Println("[helper.Delete]图片删除失败：", err)
		}
	}

	rt := rdb.Del(utils.Ctx(), fmt.Sprintf("EnterNeu:Ticket:%s", th.token)).Val()
	if rt == 1 {
		return nil
	}
	return errors.New("删除失败")
}

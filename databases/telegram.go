package databases

import (
	"encoding/json"
	"fmt"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"
	"github.com/rroy233/EnterNEU/utils"
)

type tgAllow struct {
	AdminUID  int64          `json:"adminUID"`
	AllowList map[int64]bool `json:"allowList"`
}

var tgAllowInstance *tgAllow

func GetTGAllow() *tgAllow {
	return tgAllowInstance
}

func InitTGAllow() {
	tgAllowInstance = new(tgAllow)
	cacheData := rdb.Get(utils.Ctx(), "EnterNeu:TelegramAllowList").Val()
	if cacheData == "" { //初始化
		tgAllowInstance.AdminUID = configs.Get().TGService.AdminUID
		tgAllowInstance.AllowList = make(map[int64]bool, 0)
		tgAllowInstance.save()
		return
	}
	err := json.Unmarshal([]byte(cacheData), tgAllowInstance)
	if err != nil {
		logger.FATAL.Fatalln("[databases]" + err.Error())
		return
	}
	return
}

func (ta tgAllow) save() {
	js, err := json.Marshal(ta)
	if err != nil {
		logger.Error.Println("[databases]" + err.Error())
		return
	}
	rdb.Set(utils.Ctx(), "EnterNeu:TelegramAllowList", js, -1)
}

func (ta *tgAllow) Put(uid int64) {
	if uid == ta.AdminUID {
		return
	}
	ta.AllowList[uid] = true
	ta.save()
}

func (ta *tgAllow) FlushAll() {
	ta.AdminUID = configs.Get().TGService.AdminUID
	ta.AllowList = make(map[int64]bool, 0)
	ta.save()
}

func (ta *tgAllow) Remove(uid int64) {
	ta.AllowList[uid] = false
	ta.save()
}

func (ta *tgAllow) List() string {
	out := ""
	for uid, b := range ta.AllowList {
		if b == true {
			out += fmt.Sprintf("%d\n", uid)
		}
	}
	return out
}

func (ta tgAllow) IsIn(uid int64) bool {
	if uid == ta.AdminUID {
		return true
	}
	if ta.AllowList[uid] == true {
		return true
	}
	return false
}

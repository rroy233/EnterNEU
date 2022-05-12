package databases

import (
	"fmt"
	"github.com/rroy233/EnterNEU/utils"
	"time"
)

func SaveToRedis(serviceName string, extras []string, data []byte, expTime time.Duration) error {
	keyExtra := ""
	for _, extra := range extras {
		keyExtra += ":" + extra
	}
	key := fmt.Sprintf("EnterNeu:%s%s", serviceName, keyExtra)

	err := rdb.Set(utils.Ctx(), key, data, expTime).Err()
	return err
}

func GetFromRedis(serviceName string, extras []string) (string, error) {
	keyExtra := ""
	for _, extra := range extras {
		keyExtra += ":" + extra
	}
	key := fmt.Sprintf("EnterNeu:%s%s", serviceName, keyExtra)

	data, err := rdb.Get(utils.Ctx(), key).Result()
	return data, err
}

func DeleteFromRedis(serviceName string, extras []string) error {
	keyExtra := ""
	for _, extra := range extras {
		keyExtra += ":" + extra
	}
	key := fmt.Sprintf("EnterNeu:%s%s", serviceName, keyExtra)

	err := rdb.Del(utils.Ctx(), key).Err()
	return err
}

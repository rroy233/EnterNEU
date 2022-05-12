package databases

import (
	"context"
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/logger"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitDB() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     configs.Get().Redis.Host + ":" + configs.Get().Redis.Port,
		Password: configs.Get().Redis.Password,
		DB:       configs.Get().Redis.DB,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		logger.FATAL.Fatalln("[系统服务][异常]Redis启动失败:", err)
		return
	}
	logger.Info.Println("[系统服务][成功]Redis已连接")
	return
}

package service

import (
	"github.com/rroy233/EnterNEU/configs"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"testing"
)

func TestCronClean(t *testing.T) {
	storageDIr = "../storage/"
	logger.New(false)
	configs.InitConfig("../config.yaml")
	databases.InitDB()
	cronClean()
}

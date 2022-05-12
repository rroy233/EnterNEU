package service

import (
	"fmt"
	"github.com/rroy233/EnterNEU/databases"
	"github.com/rroy233/EnterNEU/logger"
	"io/ioutil"
	"os"
	"strings"
)

func cronClean() {
	files, err := ioutil.ReadDir(storageDIr)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	for _, file := range files {
		if file.IsDir() == true || file.Name() == ".gitkeep" {
			continue
		}
		if databases.NewHelper("").SetToken(strings.Split(file.Name(), ".")[0]).Validate() == false {
			err = os.Remove(fmt.Sprintf("%s%s", storageDIr, file.Name()))
			if err != nil {
				logger.Error.Printf("[service.cronClean]文件%s删除失败\n", fmt.Sprintf("%s%s", storageDIr, file.Name()))
			} else {
				logger.Info.Printf("[service.cronClean]已成功删除文件%s\n", fmt.Sprintf("%s%s", storageDIr, file.Name()))
			}
		}
	}
	return
}
